heka-tutum
=========

This is a decoder plugin that enables decoding docker container names running in the [Tutum.co](http://tutum.co/) Docker environment.

## Status

This plugin is a work in progress and is not ready for production use. Please give it a test though and feel free to open an [issue](https://github.com/ianneub/heka-tutum/issues).

## How to use

To use this plugin you will need to add this project to your Heka source code by adding a line to `cmake/plugin_loader.cmake` that will load the plugin, like this:

    add_external_plugin(git https://github.com/ianneub/heka-tutum master)


### TutumDecoder

Example configuration:

    [DockerLogInput]
    decoder = "TutumDecoder"
    
    [TutumDecoder]
    
## Questions

Please create an issue on GitHub with any questions or comments. Pull requests are especially appreciated.

## License

See `LICENSE`.
