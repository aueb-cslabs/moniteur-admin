module.exports = {
    publicPath : process.env.NODE_ENV === 'production' ? './' : '/',
    pluginOptions: {
        electronBuilder: {
            builderOptions: {
                // options placed here will be merged with default configuration and passed to electron-builder
                "appId": "org.aueb.moniteur-admin",
                "productName": "Moniteur Admin",
                "nsis": {
                    "artifactName": "moniteur-admin-Setup-v${version}.${ext}",
                },
                "portable": {
                    "artifactName": "moniteur-admin-v${version}.${ext}",
                },
                "linux": {
                    "target": "AppImage",
                    "artifactName": "moniteur-admin-v${version}.${ext}"
                },
                "dmg": {
                    "artifactName": "moniteur-admin-Setup-v${version}.${ext}",
                },
                "extraFiles": [
                    {
                        "from": "config/config.yml",
                        "to": "config.yml"
                    }
                ],
                "publish": [
                    {
                        "provider": "github",
                        "releaseType": "release"
                    }
                ]
            }
        }
    }
}