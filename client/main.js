const { app, BrowserWindow } = require('electron');
const path = require('path');

let mainWindow;

function createWindow() {
    mainWindow = new BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            preload: path.join(__dirname, 'preload.js'),
            nodeIntegration: true,
        },
    });

    mainWindow.loadURL(
        `file://${path.join(__dirname, 'build/index.html')}`
    );

   // mainWindow.webContents.openDevTools();


    mainWindow.on('closed', () => (mainWindow = null));
}

// Secure Restorable State (Optional to prevent warning)
app.on('ready', () => {
    if (app.applicationSupportsSecureRestorableState) {
        app.applicationSupportsSecureRestorableState = () => true;
    }
    createWindow();
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') app.quit();
});

app.on('activate', () => {
    if (mainWindow === null) createWindow();
});
