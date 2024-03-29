'use strict';

import { app, protocol, BrowserWindow } from 'electron'
import { createProtocol } from 'vue-cli-plugin-electron-builder/lib'
import fs from "fs";
import yaml from 'js-yaml';
const { autoUpdater } = require("electron-updater");
const path = require('path');

const isDevelopment = process.env.NODE_ENV !== 'production';

autoUpdater.logger = require("electron-log");
autoUpdater.logger.transports.file.level = "info";

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win;

// Scheme must be registered before the app is ready
protocol.registerSchemesAsPrivileged([{scheme: 'app', privileges: { secure: true, standard: true } }]);

function createWindow () {
  // Create the browser window.
  win = new BrowserWindow({ width: 800, height: 600, autoHideMenuBar: true, webPreferences: {
    nodeIntegration: true
  } });

  if (process.env.WEBPACK_DEV_SERVER_URL) {
    // Load the url of the dev server if in development mode
    win.loadURL(process.env.WEBPACK_DEV_SERVER_URL);
    if (!process.env.IS_TEST) win.webContents.openDevTools()
  } else {
    createProtocol('app');
    // Load the index.html when not in development
    win.loadURL('app://./index.html');
    win.maximize();
  }

  win.on('closed', () => {
    win = null
  })
}

// Quit when all windows are closed.
app.on('window-all-closed', () => {
  // On macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    app.quit()
  }
});

// SSL/TSL: this is the self signed certificate support
app.on('certificate-error', (event, webContents, url, error, certificate, callback) => {
    // On certificate error we disable default behaviour (stop loading the page)
    // and we then say "it is all fine - true" to the callback
    event.preventDefault();
    callback(true);
});

app.on('activate', () => {
  // On macOS it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (win === null) {
    createWindow()
  }
});

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', async () => {
  checkAppData();
  loadConfig();
  createWindow();

  // Check for updates
  try {
    await autoUpdater.checkForUpdatesAndNotify();
  } catch (e) {
      autoUpdater.logger.log("Error occurred while calling update! " + e);
  }
});

// Exit cleanly on request from parent process in development mode.
if (isDevelopment) {
  if (process.platform === 'win32') {
    process.on('message', data => {
      if (data === 'graceful-exit') {
        app.quit()
      }
    })
  } else {
    process.on('SIGTERM', () => {
      app.quit()
    })
  }
}

function checkAppData() {
  /**
   * let execPath = app.getPath("exe") returns the
   * path of the executable file. With execPath.dir we get
   * only the directory. Adding the "/config.yml" we are able to load
   * the config in both platforms.
   */
  let execPath = path.parse(app.getPath("exe"));
  if (!fs.existsSync(app.getPath('userData')+"/config.yml")) {
    fs.openSync(app.getPath('userData')+'/config.yml', 'w');
    let finalConfig = yaml.safeLoad(fs.readFileSync(execPath.dir +'/config.yml', 'utf8'));
    fs.writeFileSync(app.getPath('userData')+'/config.yml', yaml.safeDump(finalConfig), function(err) {
      if (err) return err;
    });
  }
  else {
    let finalConfig = yaml.safeLoad(fs.readFileSync(app.getPath('userData')+'/config.yml', 'utf-8'));
    let finalConfigKeys = Object.keys(finalConfig);
    let config = yaml.safeLoad(fs.readFileSync(execPath.dir + '/config.yml', 'utf8'));
    for (let i in config) {
      if (!finalConfigKeys.includes(i.toString())) {
        finalConfig[i.toString()] = config[i.toString()];
      }
    }
    fs.writeFileSync(app.getPath('userData')+'/config.yml', yaml.safeDump(finalConfig), function(err) {
      if (err) return err;
    });
  }
}

function loadConfig() {
  try {
    global.config = yaml.safeLoad(fs.readFileSync(app.getPath('userData')+'/config.yml', 'utf8'));
    // eslint-disable-next-line no-console
  } catch (e) {
    // eslint-disable-next-line no-console
  }
}