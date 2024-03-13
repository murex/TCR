# TCR WebApp

This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 17.0.6.

## Setting up the environment

Run `npm install` to install all the dependencies.

## Development server

Run `npm start` for a dev server. Navigate to `http://localhost:4200/`. The application will automatically reload if you change any of the source files.

__Note__: For communication with TCR backend, you need to have TCR application running in development mode.
This can be done through the following command:

```bash
cd ../src
./tcr-local web -T=http
```

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

### Development

Run `npm build` to build the project for development.
The build artifacts will be stored in the `dist/` directory.

### Production

Run `npm run build-prod` to build the project for production.
The build artifacts will be stored in Go's [src/http/static/webapp](../src/http/static/webapp) directory.
This directory will then be embedded in TCR application's binary during the build process.

## Running unit tests

### One shot

Run `npm test` to execute the unit tests via [Karma](https://karma-runner.github.io).

### Watch mode

Run `npm run test-watch` to execute the unit tests in watch mode.

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via a platform of your choice. To use this command, you need to first add a package that implements end-to-end testing capabilities.

## Running linter

Run `npm run lint` to execute the linter.

## Getting the coverage report

Run `npm run coverage` to get the coverage report.

## Utility scripts

### Cleanup

Run `npm run clean` to execute the cleanup script. This script will remove all the generated files and directories.

### Reinstall dependencies

Run `npm run reinstall` to execute the reinstall script. This script will remove the `node_modules` directory and reinstall all the dependencies.

### Update dependencies

Run `npm update` to execute the update script. This script will update all the dependencies to their latest versions.

### Rebuild (development)

Run `npm run rebuild` to execute the rebuild script. This script will execute the cleanup script, reinstall script and build script in sequence.

### Rebuild (production)

Run `npm run rebuild-prod` to execute the rebuild-prod script. This script will execute the cleanup script, reinstall script and build-prod script in sequence.

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.
