[Setup]
AppName=rmbg
AppVersion=1.1
DefaultDirName={commonpf}\rmbg
DefaultGroupName=rmbg
UninstallDisplayIcon={app}\rmbg.exe
OutputBaseFilename=rmbg-setup
AppVerName=rmbg 1.1

[Files]
; Dodaj g³ówny plik aplikacji
Source: "D:\Go-Projects\rmbg\rmbg.exe"; DestDir: "{app}"; Flags: ignoreversion

[Registry]
; Dodanie opcji do kontekstowego menu dla wszystkich typów plików graficznych
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbg"; ValueType: string; ValueName: ""; ValueData: "Szybkie usuwanie t³a"; Flags: uninsdeletekey
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbg\command"; ValueType: string; ValueName: ""; ValueData: """{app}\rmbg.exe"" ""%1"""; Flags: uninsdeletekey
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbg"; ValueType: string; ValueName: "ExtendedSubCommandsKey"; ValueData: ""; Flags: uninsdeletekey
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbg"; ValueType: string; ValueName: "Position"; ValueData: "Bottom"; Flags: uninsdeletekey

Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbgm"; ValueType: string; ValueName: ""; ValueData: "Dok³adne usuwanie t³a"; Flags: uninsdeletekey
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbgm\command"; ValueType: string; ValueName: ""; ValueData: """{app}\rmbg.exe"" ""%1"" ""birefnet-massive"""; Flags: uninsdeletekey
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbgm"; ValueType: string; ValueName: "ExtendedSubCommandsKey"; ValueData: ""; Flags: uninsdeletekey
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbgm"; ValueType: string; ValueName: "Position"; ValueData: "Bottom"; Flags: uninsdeletekey
