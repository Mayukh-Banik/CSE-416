"use strict";
exports.__esModule = true;
var electron_1 = require("electron");
var path = require("path");
function createWindow() {
    var mainWindow = new electron_1.BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            preload: path.join(__dirname, '../preload/index.ts'),
            nodeIntegration: false,
            contextIsolation: true
        }
    });
    mainWindow.loadURL(process.env.ELECTRON_START_URL || "file://".concat(path.join(__dirname, '../build/index.html')));
    mainWindow.on('closed', function () {
        electron_1.app.quit();
    });
}
electron_1.app.on('ready', createWindow);
electron_1.app.on('window-all-closed', function () {
    if (process.platform !== 'darwin') {
        electron_1.app.quit();
    }
});
electron_1.app.on('activate', function () {
    if (electron_1.BrowserWindow.getAllWindows().length === 0) {
        createWindow();
    }
});
