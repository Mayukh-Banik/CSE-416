"use strict";
exports.__esModule = true;
var electron_1 = require("electron");
electron_1.contextBridge.exposeInMainWorld('electron', {
    send: function (channel, data) {
        electron_1.ipcRenderer.send(channel, data);
    },
    receive: function (channel, func) {
        electron_1.ipcRenderer.on(channel, function (event) {
            var args = [];
            for (var _i = 1; _i < arguments.length; _i++) {
                args[_i - 1] = arguments[_i];
            }
            return func.apply(void 0, args);
        });
    }
});
