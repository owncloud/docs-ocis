/*
 * To be parsed via: https://structurizr.com/dsl
 * (Help available via selector on the left side below "Upload"
 * Then exported via button PlantUML
 * Then imported at http://www.plantuml.com/plantuml
 * Then exported as .svg image
*/

workspace "ownCloud Infinite Scale" "The oCIS C4 Model" {

	model {
		user = person "Client" "Clients accessing oCIS"
		ocis = softwareSystem "oCIS environment" "ownCloud Infinite Scale"
		external_storage = softwareSystem "External Storage" "Google, FTP, Cloud ..." "Existing System"
		internal_storage = softwareSystem "Internal Backend Storage" "NFS, SMB, CephFS" "Existing System"
		user -> ocis "Connects"
		ocis -> internal_storage "API" "Uses"
		ocis -> external_storage "API" "Uses"
	}

	views {
		systemlandscape "SystemLandscape" {
			include *
			autoLayout
		}

		systemContext ocis "SystemContext" "An example of a System Context diagram." {
			include *
			autoLayout
		}

		styles {
			element "Software System" {
				background #1168bd
				color #ffffff
			}
			element "Existing System" {
				background #999999
				color #ffffff
			}
			element "Person" {
				shape person
				background #08427b
				color #ffffff
			}
		}
	}
}
