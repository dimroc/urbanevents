'use strict';

var gulp           = require('gulp');
var del            = require('del');
var mainBowerFiles = require('main-bower-files')
var concat         = require('gulp-concat')

// Load plugins
var $            = require('gulp-load-plugins')();
var browserify   = require('browserify');
var watchify     = require('watchify');
var source       = require('vinyl-source-stream');
var sourceFile   = './app/scripts/app.js';
var destFolder   = './dist/scripts';
var destFileName = 'app.js';
var browserSync  = require('browser-sync');
var reload       = browserSync.reload;

// Styles
gulp.task('styles', ['sass']);

gulp.task('sass', function() {
    return gulp.src(['app/styles/**/*.scss', 'app/styles/**/*.css'])
        .pipe($.rubySass({
            style: 'expanded',
            precision: 10,
            loadPath: ['app/bower_components']
        }))
        .pipe($.autoprefixer('last 1 version'))
        .pipe(gulp.dest('dist/styles'))
        .pipe($.size());
});

var bundler = watchify(browserify({
    entries: [sourceFile],
    debug: true,
    insertGlobals: true,
    cache: {},
    packageCache: {},
    fullPaths: true
}));

bundler.on('update', rebundle);
bundler.on('log', $.util.log);

function rebundle() {
    return bundler.bundle()
        // log errors if they happen
        .on('error', $.util.log.bind($.util, 'Browserify Error'))
        .pipe(source(destFileName))
        .pipe($.streamify($.replace(/__ENV_CITYSERVICE_URL__/g, '"' + process.env.CITYSERVICE_URL + '"')))
        .pipe(gulp.dest(destFolder))
        .on('end', function() {
            reload();
        });
}

// Scripts
gulp.task('scripts', rebundle);

gulp.task('buildScripts', function() {
    return browserify(sourceFile)
        .bundle()
        .pipe(source(destFileName))
        .pipe($.streamify($.replace(/__ENV_CITYSERVICE_URL__/g, '"' + process.env.CITYSERVICE_URL + '"')))
        .pipe(gulp.dest('dist/scripts'));
});

// HTML
gulp.task('html', function() {
    return gulp.src('app/*.html')
        .pipe($.useref())
        .pipe(gulp.dest('dist'))
        .pipe($.size());
});

// Images
gulp.task('images', function() {
    return gulp.src('app/images/**/*')
        .pipe($.cache($.imagemin({
            optimizationLevel: 3,
            progressive: true,
            interlaced: true
        })))
        .pipe(gulp.dest('dist/images'))
        .pipe($.size());
});

// Fonts
gulp.task('fonts', function() {
            filter: '**/*.{eot,svg,ttf,woff,woff2}'
    return gulp.src(mainBowerFiles({
        }).concat('app/fonts/**/*'))
        .pipe(gulp.dest('dist/fonts'));
});

// Clean
gulp.task('clean', function(cb) {
    $.cache.clearAll();
    cb(del.sync(['dist/styles', 'dist/scripts', 'dist/images']));
});

// Bundle
gulp.task('bundle', ['styles', 'scripts', 'bower', 'bowerCss'], function() {
    return gulp.src('./app/*.html')
        .pipe($.useref.assets())
        .pipe($.useref.restore())
        .pipe($.useref())
        .pipe(gulp.dest('dist'));
});

gulp.task('buildBundle', ['styles', 'buildScripts', 'bower', 'bowerCss'], function() {
    return gulp.src('./app/*.html')
        .pipe($.useref.assets())
        .pipe($.useref.restore())
        .pipe($.useref())
        .pipe(gulp.dest('dist'));
});

// Bower helper
gulp.task('bower', function() {
    gulp.src('app/bower_components/**/*.js', {
            base: 'app/bower_components'
        })
        .pipe(gulp.dest('dist/bower_components/'));

});

//
// concat just leaflet/*.css to `vendor.css`
//
gulp.task('bowerCss', function() {
  var cssFilter = (/app\/bower_components\/leaflet\/.*css/i);
  return gulp.src(mainBowerFiles({filter: cssFilter}))
    .pipe(concat('vendor.css'))
    .pipe(gulp.dest('dist/styles'))
})

gulp.task('json', function() {
    gulp.src('app/scripts/json/**/*.json', {
            base: 'app/scripts'
        })
        .pipe(gulp.dest('dist/scripts/'));
});

// Robots.txt and favicon.ico
gulp.task('extras', function() {
    return gulp.src(['app/*.txt', 'app/*.ico'])
        .pipe(gulp.dest('dist/'))
        .pipe($.size());
});

// Watch
gulp.task('watch', ['html', 'fonts', 'bundle'], function() {

    browserSync({
        notify: false,
        logPrefix: 'BS',
        // Run as an https by uncommenting 'https: true'
        // Note: this uses an unsigned certificate which on first access
        //       will present a certificate warning in the browser.
        // https: true,
        server: ['dist', 'app']
    });

    // Watch .json files
    gulp.watch('app/scripts/**/*.json', ['json']);

    // Watch .html files
    gulp.watch('app/*.html', ['html']);

    gulp.watch(['app/styles/**/*.scss', 'app/styles/**/*.css'], ['styles', reload]);



    // Watch image files
    gulp.watch('app/images/**/*', reload);
});

// Build
gulp.task('build', ['html', 'buildBundle', 'images', 'fonts', 'extras'], function() {
    gulp.src('dist/scripts/app.js')
        .pipe($.uglify())
        .pipe($.stripDebug())
        .pipe(gulp.dest('dist/scripts'));
});

gulp.task('clearCityserviceUrlEnv', function() {
  process.env.CITYSERVICE_URL = "";
});

gulp.task('copyToService', ['clearCityserviceUrlEnv', 'build'], function() {
  gulp.src('dist/**/*')
    .pipe(gulp.dest('../cityservice/public'));
});

gulp.task('deploy', ['clean', 'copyToService']);

// Default task
gulp.task('default', ['clean', 'build' ]);
