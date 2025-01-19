'use strict'

/**
 * Registers the detect-unused-images extension.
 *
 * @param {Object} config - The configuration object.
 * @param {Array<string>} [config.excludeimageextension] - List of image extensions to exclude from detection.
 */
module.exports.register = function ({config}) {
    const logger = this.getLogger('detect-unused-objects')
    logger.info('Start unused objects detection')

    // If config.excludeExtension is provided, merge it with the default extensions to ignore
    const extensionToIgnore = []
    if(config.excludeextension){
        config.excludeextension.forEach(element => extensionToIgnore.push(element))
        // extensionToIgnore.forEach(element => console.log(element))
    }

        logger.info('Assets with extensions will be ignored %s', Array.from(extensionToIgnore))

    this.on('contentClassified', ({ contentCatalog }) => {

        // regex to identify if there are keys in pages that contain includable objects
        // note that audio and video family directories are currenly not implemented antora
        // note only use macros that have a family directory
        // g1 ... any character before g2, allowed: whitespaces, start of string and .macro at the line start
        // g2 ... any of the named macros, note that the : is a mandatory part
        // g3 ... optional for an additional : (like image can have : and ::)
        // g4 ... the relevant path, only use if it not empty
        const rg = /([\s*]|^[\.]|\A)(image:|include:|xref:)(:{1}|^\r)?([^\s|^\[]*)/g

        // array to define which macro stores/refernces its files in which default family root folder
        // note that include and xref can reference, if not otherwise defined, also to the
        // attachments and examples folder
        // key = part of the path like partials$, value = default family directory if not otherwise defined
        const fam_dir_arr = {attachment: 'attachments',
                             example: 'examples',
                             include: 'pages',
                             partial: 'partials',
                             image: 'images',
                             xref: 'pages'
                             //audio: 'images',      // family currently not implemented in Antora
                             //video: 'images'       // family currently not implemented in Antora
                            }

        contentCatalog.getComponents().forEach(({ versions }) => {
          versions.forEach(({ name: component, version }) => {

            const allFiles = get_all_files (contentCatalog, extensionToIgnore)
            // allFiles.forEach(element => console.log(element))

            const dir_sep = get_dir_seperator(allFiles)

            const pages = contentCatalog.findBy({ component, version, family: 'page' })
            const pageReferences = get_path_objects (pages, rg, fam_dir_arr, dir_sep, this)
            // pageReferences.forEach(element => console.log(element))

            const partials = contentCatalog.findBy({ component, version, family: 'partial' })
            const partialReferences = get_path_objects (partials, rg, fam_dir_arr, dir_sep, this)
            // partialReferences.forEach(element => console.log(element))

            const nav = contentCatalog.findBy({ component, version, family: 'nav' })
            const navReferences = get_path_objects (nav, rg, fam_dir_arr, dir_sep, this)
            // navReferences.forEach(element => console.log(element))

            const navFiles = get_all_nav_files (nav)
            // navFiles.forEach(element => console.log(element))

            // collect all references into an array to sort before adding to the set
            const coll = []
            pageReferences.forEach((element) => coll.push(element))
            partialReferences.forEach((element) => coll.push(element))
            navReferences.forEach((element) => coll.push(element))
            const allUniqueReferences = new Set(coll.sort())
            // allUniqueReferences.forEach(element => console.log(element))

            // remove found objects navigation files from allFiles list
            const allFilesNoNav = allFiles.filter(function (x) {
               const not_found = !navFiles.has(x)     // needs to be a set
               // console.log(not_found, x)
               return not_found
            })
            // remove found reference objects from allFilesNoNav list
            const orphandArray = allFilesNoNav.filter(function (x) {
               const not_found = !allUniqueReferences.has(x)     // needs to be a set because of has
               // console.log(not_found, x)
               return not_found
            }).filter(Boolean)                                   // removy any undefined elements

//console.log('aaaa')
console.log(orphandArray)

this.stop()
        })
      })
//this.stop()

/*
                }catch(error){
                    logger.fatal('%s (%s) - %s',file.src.component, file.src.version, file.src.basename)
                    logger.fatal(error);
                }
        });
//console.log(objectReferences)
        logger.info('Found %s images references', objectReferences.size)
*/
/*
        contentCatalog.getFiles()
            .filter(file => file.src.family === 'image' && !extensionToIgnore.has(file.src.extname))
            .forEach(img => {
                // check if the media are used in the module or in an external
                // ex: ROOT:images/myImage.png or images:myImages.png
                if(!(objectReferences.has(img.src.relative.toString()) || objectReferences.has(img.src.module + ':' + img.src.relative.toString()))){
                    unusedObjects.add(img)
                    logger.warn('[%s] [%s] %s', img.src.component, img.src.version,img.src.path)
                }
            })
        logger.info('Finish and detecting %s unused images', unusedObjects.size);

        if( unusedObjects.size > 0 ){
            logger.warn('Some images are unused, check previous logs and delete unused images.')
        }
*/
    })
}

function get_all_files (contentCatalog, extensionToIgnore) {
    // get a list of all files from the contentCatalog
    // no need to have a set, as files are unique
    // exclude files that match an extension from an config array list

    const files = []

    if (extensionToIgnore.length === 0) {    // get all files
      contentCatalog.getFiles()
        .forEach(file => {
           const item = '' + file.src.path
           files.push(item)
      })
    } else {                                     // only files that dont match an extension pattern
      contentCatalog.getFiles()
        .forEach(file => {
           const item = '' + file.src.path       // this can cause string 'undefined'
           let found
           found = extensionToIgnore.some(ext => {
             return item.endsWith(ext)           // return and exit some if true
           })
           if (!found && item != 'undefined') {  // only if found and not the string 'undefined'
             files.push(item)
           }
         })
    }

    return  files.sort().filter(Boolean)         // return, sort, remove true undefined
}

function get_all_nav_files (navigationCatalog) {
    // get a list of all navigation files

    const files = new Set
    for (const item of navigationCatalog) {
      files.add(item.src.path)
    }
    return files
}

function get_dir_seperator(allFiles) {
    // identify and return the directory separator used by the OS
    // this is needed when assembling the path in get_path_objects

    let dir_sep
    if (allFiles.length) {
      dir_sep = allFiles[0].includes('/') ? '/' : '\\' // either linux '/' or windows '\\'
    } else {
      dir_sep = '/' // defaults to linux based directory separator if array is empty to avoid errors
    }
    return dir_sep
}

function get_path_objects (pages, rg, fam_dir_arr, dir_sep, x) {
    // get all paths that are defined in macros via a regex
    // the array returned contains components to easy reassemble a path to compare
    // x: is the 'this' object and only needed for development like when using x.stop()

    // the final path assembled looks like the following, which matches the scheme of get_all_files()
    // 'modules/<module>/<family>/path/file.ext'
    // like 'modules/ROOT/pages/abc/file.adoc'

    const s1 = []
    const s2 = []
    const objectReferences = []

    // regex over all pages, only group 2 and 4 are relevant
    for (const page of pages) {
      const iterator  = [...page._contents.toString().matchAll(rg)]
      const g2 = iterator.map(m => m[2]) // macro
      const g4 = iterator.map(m => m[4]) // path

      if (g2.length > 0) {                        // only if at least one match was found in a page
        for (let i=0; i < g2.length; i++) {       // for all matches found do

          if (g2[i] && g4[i]) {                   // only if there is content in groups 2 and 4
            const macro = g2[i].replace(/:$/, '') // remove any colon if exists, needed to get array value

            let [family, path] = g4[i].split('$') // check if we have a family coordinate like partials$
            if (path === undefined || path.length == 0) {  // path is empty by if no $
              [family, path] = [path, family]     // if path is empty = no family coordinate, swap the vars
              family = macro                      // populate the default familiy coordinate from the array

            }

            if (path.includes('{')) {             // if path contains a '{' == attribute usage 
              continue                            // refuse, we cant resolve attributes to get the final path
            }

            if (path.includes('#')) {             // if path contains a '#' == xref section ref usage 
              continue                            // refuse, we do not check section references
            }

            const isHttp = path.toLowerCase().startsWith('http')
            if (isHttp) {                         // macros like include and xref can contain an URL
              continue                            // refuse, we do not check external resources
            }

            if (path.includes('@')) {             // if path contains a '@' == xref version reference 
              continue                            // refuse, we do not check version references
            }                                     // version@component, remaining: module = ok

            const a = path.split(dir_sep).pop()   // get the last path component - if it is a path
            const isFile = a.search(/([?.+].+)/g) // if there is a . it is a file, else it is a xref refernce
            if (isFile < 0) {                     // skip xref section references
              continue                            // refuse, section references are no files
            }

            if (path.startsWith('./') || path.startsWith('.\\')) { // reconstruct the full path component if relative
              const basename = page.src.basename  // the pure file name
              const relative = page.src.relative  // the relative directory of the page
              const re_wo_bas = relative.substring(0,relative.lastIndexOf(basename)-1)
              path = re_wo_bas + path.substring(1) // assemble a correct FULL relative path
            }

            path = path.replaceAll('//', '/')     // linux only: replace any path occurrences of '//' with '/'
                                                  // windows   : does not allow double path separators '\\' = np

            if (path.startsWith('/') || path.startsWith('\\')) {   // if path start with '/' or '\'
              path = path.substring(1)            // remove first path separator if extists
            }

            if (path.includes(':')) {                 // path contains a module component like ROOT:pages/xyz/...
              path = path.replaceAll(':', dir_sep)    // make it a real path like ROOT/pages/xyz/file.adoc
            } else {                                  // manually add module and family
              path = fam_dir_arr[family] + dir_sep + path  // xyz/file.adoc --> pages/xyz/file.adoc
              path = page.src.module + dir_sep + path // pages/xyz/file.adoc --> ROOT/pages/xyz/file.adoc
            }                                         // now we have a guaranteed complete path 

            path = 'modules' + dir_sep + path         // ROOT/pages/xyz/file.adoc --> modules/ROOT/pages/xyz/file.adoc

            s1.push(macro, family, path)          // technically only path and family is necessary
//if (path.includes('systemd.adoc')) {
//console.log(g2[i], g4[i])
//console.log(family, path)
//x.stop()
//}

          }                                       // we use all three for debuggiung purposes
        }
      }
    }

    // make parsable blocks for each object
    s1.forEach((element) => s2.push(element))

    // note there must be the same number of splice elements as defined in s1.push above 
    while(s2.length) objectReferences.push(s2.splice(0,3))

    // now we have the correct raw data
    // create the final set of objects = uniqe entries
    const uniqueObjectRefernces = new Set()

    // the path is on the third location
    for (let i=0; i < objectReferences.length; i++) {
      uniqueObjectRefernces.add(objectReferences[i][2])
    }

    return uniqueObjectRefernces
}
