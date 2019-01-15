var ManifestPlugin = require('webpack-manifest-plugin');

module.exports = {
    baseUrl: '/',
    configureWebpack: {
        plugins: [
            new ManifestPlugin()
        ]
    }
}
