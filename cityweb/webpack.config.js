var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require("extract-text-webpack-plugin");

var stylusLoader = ExtractTextPlugin.extract(
  'style-loader',
  'css-loader?module&localIdentName=[name]__[local]___[hash:base64:5]' +
    '&disableStructuralMinification' +
  '!autoprefixer-loader!' +
  'stylus-loader?paths=src/app/client/css/&import=./ctx'
);

var plugins = [
    new webpack.NoErrorsPlugin(),
    new webpack.optimize.DedupePlugin(),
    new ExtractTextPlugin('bundle.css'),
];

if (process.env.NODE_ENV === 'production') {
  plugins = plugins.concat([
    new webpack.optimize.UglifyJsPlugin({
      output: {comments: false},
      test: /bundle\.js?$/
    }),
    new webpack.DefinePlugin({
      'process.env': {NODE_ENV: JSON.stringify('production')}
    })
  ]);
  var stylusLoader = ExtractTextPlugin.extract(
    'style-loader',
    'css-loader?module&disableStructuralMinification' +
      '!autoprefixer-loader' +
      '!stylus-loader?paths=src/app/client/css/&import=./ctx'
  );
};

var sassLoader = ExtractTextPlugin.extract(
  'style-loader',
  'css-loader?module&disableStructuralMinification' +
  '!autoprefixer-loader' +
  '!sass-loader?includePaths[]=' + path.resolve(__dirname, './src')
);

var config  = {
  entry: {
    bundle: path.join(__dirname, 'src/app/client/index.js')
  },
  output: {
    path: path.join(__dirname, 'src/app/server/data/static/build'),
    publicPath: "/static/build/",
    filename: '[name].js'
  },
  plugins: plugins,
  module: {
    loaders: [
      {test: /\.scss$/, loader: sassLoader},
      {test: /\.styl$/, loader: stylusLoader},
      {test: /\.(png|gif)$/, loader: 'url-loader?name=[name]@[hash].[ext]&limit=5000'},
      {test: /\.svg$/, loader: 'url-loader?name=[name]@[hash].[ext]&limit=5000!svgo-loader?useConfig=svgo1'},
      {test: /\.(pdf|ico|jpg|eot|otf|woff|ttf|mp4|webm)$/, loader: 'file-loader?name=[name]@[hash].[ext]'},
      {test: /\.json$/, loader: 'json-loader'},
      {test: /\.js$/, loader:'imports?jQuery=jquery,$=jquery'},
      {
        test: /\.jsx?$/,
        include: path.join(__dirname, 'src/app/client'),
        loaders: ['babel']
      },

      // -- Bootstrap Sass --
      // **IMPORTANT** This is needed so that each bootstrap js file required by
      // bootstrap-sass-loader has access to the jQuery object
      { test: /bootstrap-sass\/assets\/javascripts\//, loader: 'imports?jQuery=jquery' },
      //{ test: /\.scss$/, loader: "style!css!sass?outputStyle=expanded" },

      // ToDo: custom path and source map option did not work
      //{ test: /\.scss$/,
      //  loader: "style!css!sass?outputStyle=expanded&sourceMap=true&includePaths[]=" + bootstrapPathStylesheets },

      // Needed for the css-loader when [bootstrap-sass-loader](https://github.com/justin808/bootstrap-sass-loader)
      // loads bootstrap's css.
      { test: /\.woff(\?v=\d+\.\d+\.\d+)?$/,   loader: "url?limit=10000&minetype=application/font-woff" },
      { test: /\.woff2(\?v=\d+\.\d+\.\d+)?$/,  loader: "url?limit=10000&minetype=application/font-woff" },
      { test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,    loader: "url?limit=10000&minetype=application/octet-stream" },
      { test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,    loader: "file" },
      { test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,    loader: "url?limit=10000&minetype=image/svg+xml" }
    ]
  },
  resolve: {
    extensions: ['', '.js', '.jsx', '.styl', '.scss'],
    alias: {
      '#app': path.join(__dirname, '/src/app/client'),
      '#c': path.join(__dirname, '/src/app/client/components'),
      '#css': path.join(__dirname, '/src/app/client/css')
    }
  },
  svgo1: {
    multipass: true,
    plugins: [
      // by default enabled
      {mergePaths: false},
      {convertTransform: false},
      {convertShapeToPath: false},
      {cleanupIDs: false},
      {collapseGroups: false},
      {transformsWithOnePath: false},
      {cleanupNumericValues: false},
      {convertPathData: false},
      {moveGroupAttrsToElems: false},
      // by default disabled
      {removeTitle: true},
      {removeDesc: true}
    ]
  }
};

module.exports = config;
