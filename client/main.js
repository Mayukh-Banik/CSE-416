const { app, BrowserWindow, ipcMain, dialog} = require('electron');
const path = require('path');
const fs = require('fs')

let mainWindow;

function createWindow() {
    mainWindow = new BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            preload: path.join(__dirname, 'preload.js'),
            contextIsolation: true,
            nodeIntegration: false,
            enableRemoteModule: false,
        },
    });

    mainWindow.loadURL(
        `file://${path.join(__dirname, 'build/index.html')}`
    );

    // CLOSE LATER
   mainWindow.webContents.openDevTools();
    // CLOSE LATER

    mainWindow.on('closed', () => (mainWindow = null));
}

ipcMain.handle('save-file',async (event,{fileName, fileData})=>{
    try {
        const directoryPath = path.join(__dirname, "..", 'squidcoinFiles'); // Adjust as needed for your file path
        const filePath = path.join(directoryPath, fileName);
        
        // Ensure the directory exists
        if (!fs.existsSync(directoryPath)) {
            fs.mkdirSync(directoryPath, { recursive: true }); // Create the directory if it doesn't exist
        }

        fs.writeFileSync(filePath, fileData);
        return {success: true, message: 'File saved successfully'};
        
    } catch(error){
        console.log('Error saving file: ', error);
        return {success: false, message: "File was"}
    }

})

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

