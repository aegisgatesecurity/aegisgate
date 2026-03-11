const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('aegisgate', {
  onAppReady: (callback) => ipcRenderer.on('app-ready', callback),
  sendCheckStatus: () => ipcRenderer.send('check-status'),
  exit: () => ipcRenderer.send('exit')
});