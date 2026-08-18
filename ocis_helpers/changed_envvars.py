import yaml
import sys
import os
from datetime import date
from urllib.request import urlopen

## this python script generates adoc files for added, deprecated and removed
## envvars based on the 'env_vars.yaml' that must exist in each referenced version.
## it is CRUCIAL that versions compared are actual - do required updates first!
#
## note the helpers have been migrated in summer 2026 from the ocis repo to docs.
## only 2 stable (production) versions are kept which you can compare. 
## note that we are always comparing from github sources and NOT local files

## when the files got created, you MUST do some post work manually like referencing services with xref:
## when executing, files get recreated, existing content will be overwritten!!

## !! you MUST run this script from the ocis_helper folder !!
## like: python3 changed_envvars.py
#
## create a branch beforehand to prepare for the changes

## WHAT TO CHANGE FOR A NEW RELEASE
##
## versionOld:     This is the version you are comparing FROM
## versionNew:     This is the version you are comparing TO
## excludePattern: add the FROM release to the list. The code then searches for anything that does not match the exclude pattern.

# CHANGE according your needs (version comparison):
# 'versionOld' is the base version to compare from
# 'versionNew' is the target version to compare to
# branches are written the follwoing: 'master' or 'x,y' such as '8.1'
versionOld: str    = '8.1'
versionNew: str    = '8.2'

# CHANGE according your needs (for printing, to_version also to identify all other findings that should be excluded)
from_version: str  = '8.1.0'
to_version: str    = '8.2.0'

# CHANGE according your needs (for file names to be created)
# this will create files such as 8.0.0-8.1.0-added, 8.0.0-8.1.0-removed etc.
# this must match which versions you compare.
nameComponent: str = '8.1.0-8.2.0'

# ADD elements of version that have been published so they get EXCLUDED from gathering.
# from version 7 onwards we dont use patch versions anymore - they should be present.
# fix any patch version in the ocis repo first and redo the whole helper process.
# the last excluded entry should be the version you take as comparison base which is 'from_version'.
defaultExcludePattern: list[str] = ['pre5.0', '5.0',
									'6.0', '6.0.0', '6.0.1', '6.1.0', '6.7',
									'7.0', '7.0.0', '7.1.0', '7.2.0', '7.3.0',
									'8.0.0', '8.1.0']

# IMPORTANT: there may be additional unexpected versions that may appear from backporting in ocis which
# must be excluded too. if found, they are printed to the console for furter investigation.
extraExcludePattern: list[str]   = ['8.3.0']

# DO NOT CHANGE
# this is the sub-path the added/deprecated and removed files are written to
# this path is independent of other helpers, although the first patch component (services) must match
adocWritePath: str = 'services/env_var_deltas/'

# DO NOT CHANGE
# this is the sub-path the 'env_vars.yaml' file is located to read
# this path MUST align with the other helpers 
adocReadPath: str = 'services/persistent_files/'

# DO NOT CHANGE
# define the path elements for github raw access
# the right component is identical to local relative but will be updated by the code for github raw access
git_left_dir: str  = 'https://raw.githubusercontent.com/owncloud/docs-ocis/refs/heads/'
git_right_dir: str = adocReadPath + 'env_vars.yaml'

##
## functions to be called by main
##

def merge_exclude_lists(defaultExcludePattern, extraExcludePattern) -> list[str]:
	# merge the arrays and remove duplicates if any
	merged_array = list(set(defaultExcludePattern) | set(extraExcludePattern))
	merged_array.sort()
	return merged_array

def get_cli_dryrun_param() -> bool:
	# exit if there is no parameter added
	if (len(sys.argv) -1) == 0:
		return False
	
	# we only allow one cli parameter which is -h (help) or -d (dryrun)
	param_1 = sys.argv[1]
	match param_1:
		case '-d':
			return True
		case '-h' | _:
			print('Allowed CLI parameters are: -h (this help message ) or -d (dryrun), exiting\n')
			sys.exit()

def correct_git_right_dir() -> None:
	# get the foldername of the current directory and correct the right git side fo rgithub raw access
	# by that we are independent of any folder names or levels the helpers are located
	global git_right_dir

	last_segment = os.path.basename(os.getcwd()) + '/'
	git_right_dir = last_segment + git_right_dir

def create_target_directory() -> None:
	# create the target folders if not exists
	if not os.path.exists(adocWritePath):
		try:
			os.makedirs(adocWritePath, exist_ok=True)
		except OSError:
			if not os.path.isdir(adocWritePath):
				raise
				sys.exit()

def get_sources(versionOld: str, versionNew: str) -> (str, str):
	urlOld = git_left_dir + versionOld + '/' + git_right_dir
	urlNew = git_left_dir + versionNew + '/' + git_right_dir

	try:
		fileOld = urlopen(urlOld).read().decode('utf-8')
		fileNew = urlopen(urlNew).read().decode('utf-8')
		print('Reading the following files for comparison:\n')
		print('  ' + urlOld)
		print('  ' + urlNew)
		return	yaml.safe_load(fileOld), yaml.safe_load(fileNew)

	except Exception as e:
		print(e)
		print('  ' + urlOld)
		print('  ' + urlNew)
		sys.exit()

def get_added(fileNew: str, excludePattern: list[str], to_Version: str) -> set[str]:
	# create dict with added items
	addedWith: set[str] = {}
	additionally_addedWith: set[str] = {}
	for key, value in fileNew.items():
		if not fileNew[key]['introductionVersion'] in str(excludePattern):
			# the requirement says that we only want to catch the target version
			if fileNew[key]['introductionVersion'] == to_Version:
				addedWith[key] = value
			# but if we find anonter version, we should add it to another list which later can be printed
			else:
				additionally_addedWith[key] = fileNew[key]['introductionVersion']
	return addedWith, additionally_addedWith

def get_removed(fileOld: str, fileNew: str) -> set[str]:
	# create dict with removed items
	removedWith: set[str] = {}
	for key, value in fileOld.items():
		if not key in fileNew:
			removedWith[key] = value
	return removedWith

def get_deprecated(fileNew: str) -> set[str]:
	# create dict with deprecated items
	deprecatedWith: set[str] = {}
	for key, value in fileNew.items():
		if value['removalVersion']:
			deprecatedWith[key] = value
	return deprecatedWith

def create_adoc_start(type_text: str, from_version: str, to_version: str, creation_date: str, columns: str, closing: str) -> str:
	# create the page/table header
	# 'closing' contains variable column names dependen if added/removed ir deprecated
	a: str = '''// # {ftype} Variables between oCIS {ffrom} and oCIS {fto}
// commenting the headline to make it better includable

// table created per {fdate}
// the table should be recreated/updated on source () changes

[width="100%",cols="{fcolumns}",options="header"]
|===
| Service | Variable | Description | {fclosing}

'''.format(ftype = type_text, ffrom = from_version, fto = to_version, fdate = creation_date, fcolumns = columns, fclosing = closing)
	return a

def create_adoc_end() -> str:
	# close the table
	a: str = '''|===

'''
	return a

def add_adoc_line_1(service: str, variable: str, description: str, value: str) -> str:
	# add a table line for added/removed
	# the dummy values are only here to have the same number of parameters as add_adoc_line_2
	a: str = '''| {fservice}
| {fvariable}
| {fdescription}
| {fvalue}

'''.format(fservice = service, fvariable = variable, fdescription = description, fvalue = value)
	return a

def add_adoc_line_2(service, variable, description, removalVersion, deprecationInfo):
	# add a table line for deprecated, this has different columns
	a: str = '''| {fservice}
| {fvariable}
| {fdescription}
| {fremovalVersion}
| {fdeprecationInfo}

'''.format(fservice = service, fvariable = variable, fdescription = description, fremovalVersion = removalVersion, fdeprecationInfo = deprecationInfo)
	return a

def create_table(type_text: str, source_dict: set[str], from_version: str, to_version: str, date_today: str, type: bool = False) -> str:
	# get the table header
	columns: str = '~,~,~,~' if type == False else '~,~,~,~,~'
	closing: str = 'Default' if type == False else 'Removal Version | Deprecation Info'
	a: str = create_adoc_start(type_text, from_version, to_version, date_today, columns, closing)

	if not type:
	# added and removed envvars
		# first add all ocis_
		for key, value in source_dict.items():
			if key.startswith('OCIS_'):
				a += add_adoc_line_1(
						# note that any envvar starting with ocis cant be assigned to a service automatically
						# the xref must be corrected in the output file manually
						'xref:deployment/services/env-vars-special-scope.adoc[Special Scope Envvars]',
						key,
						value['description'],
						value['defaultValue']
					)
		# then add all others
		for key, value in source_dict.items():
			if not key.startswith('OCIS_'):
				a += add_adoc_line_1(
						'xref:{s-path}/xxx.adoc[xxx]',
						key,
						value['description'],
						value['defaultValue']
					)
	else:
	# deprecated envvars
		# first add all ocis_
		for key, value in source_dict.items():
			if key.startswith('OCIS_'):
				a += add_adoc_line_2(
						'xref:deployment/services/env-vars-special-scope.adoc[Special Scope Envvars]',
						key,
						value['description'],
						value['removalVersion'],
						value['deprecationInfo']
					)
		# then add all others
		for key, value in source_dict.items():
			if not key.startswith('OCIS_'):
				a += add_adoc_line_2(
						'xref:{s-path}/xxx.adoc[xxx]',
						key,
						value['description'],
						value['removalVersion'],
						value['deprecationInfo']
					)

	# finally close the table
	a += create_adoc_end()
	return a

def write_output(a: str, type_text: str) -> None:
	# write the content to a file
	try:
		with open(adocWritePath + '/' + nameComponent + '-' + type_text + '.adoc', 'w') as file:
			file.write(a)
	except Exception as e:
		print('Failed creating ' + type_text + ' file')
		print(e)
		sys.exit()

def main() -> None:
	## here are the tasks in sequence

	dry_run: bool						= False
	fileOld: str						= ''
	fileNew: str						= ''
	addedWith: set[str]					= {}
	additionally_addedWith: set[str]	= {}
	removedWith: set[str]				= {}
	deprecatedWith: set[str]			= {}
	excludePattern: list[str]			= []

	dry_run = get_cli_dryrun_param()

	correct_git_right_dir()
	create_target_directory()
	excludePattern						= merge_exclude_lists(defaultExcludePattern, extraExcludePattern)
	fileOld, fileNew					= get_sources(versionOld, versionNew)
	addedWith, additionally_addedWith	= get_added(fileNew, excludePattern, to_version)
	removedWith							= get_removed(fileOld, fileNew)
	deprecatedWith						= get_deprecated(fileNew)

	# additional, unexpected introduction versions have been found.
	# although we just could exclude them automatically,
	# we print them to investigate and for possible furter exclusion in the variable above.
	if bool(additionally_addedWith):
		print('\nThis additional out-of-scope envvars have been found in %s. They may require fixing in ocis or extra excluded here:\n' % to_version)
		for key, value in additionally_addedWith.items():
			print(key, value)
		print('')

	a: str = create_table('Added',      addedWith,      from_version, to_version, date.today().strftime('%Y.%m.%d'))
	r: str = create_table('Removed',    removedWith,    from_version, to_version, date.today().strftime('%Y.%m.%d'))
	d: str = create_table('Deprecated', deprecatedWith, from_version, to_version, date.today().strftime('%Y.%m.%d'), True)

	if dry_run:
		print('Creation of tables succeeded not being written due to the dryrun flag set\n')
		sys.exit()

	write_output(a, 'added')
	write_output(r, 'removed')
	write_output(d, 'deprecated')

	print('\nSuccess, see files created in: ' + adocWritePath)
	print('First: check and update the files manually to group entries for easier reading and to fix service names (xxx).')
	print('Second: you must run n Antora build. The files written are not tracked. Only a build makes them trackable.')

if __name__ == '__main__':
	main()
