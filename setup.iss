[Setup]
AppName=RMBG
AppVersion=1.2
VersionInfoVersion=1.2.0.0
DefaultDirName={commonpf}\RMBG
DefaultGroupName=RMBG
DisableProgramGroupPage=yes
UninstallDisplayIcon={app}\RMBG.exe
OutputDir=build\bin
OutputBaseFilename=RMBG-setup
AppVerName=RMBG 1.2

[Languages]
Name: "en"; MessagesFile: "compiler:Default.isl"
Name: "pl"; MessagesFile: "compiler:Languages\Polish.isl"

[CustomMessages]
en.AddContextOption=Add a background removal option to the context menu
pl.AddContextOption=Dodaj opcję usuwania tła do menu kontekstowego

en.ContextQuick=Quick background removal
pl.ContextQuick=Szybkie usuwanie tła

en.ContextPrecise=Precise background removal
pl.ContextPrecise=Dokładne usuwanie tła

en.AdditionalIcons=Additional icons
pl.AdditionalIcons=Dodatkowe ikony

en.AdditionalOptions=Additional options
pl.AdditionalOptions=Dodatkowe opcje

en.CreateDesktopIcon=Create a desktop icon
pl.CreateDesktopIcon=Utwórz ikonę na pulpicie

en.CreateStartMenuIcon=Create a Start Menu icon
pl.CreateStartMenuIcon=Utwórz ikonę w menu Start

[Icons]
Name: "{group}\RMBG"; Filename: "{app}\RMBG.exe"; Tasks: "startmenuicon"
Name: "{commondesktop}\RMBG"; Filename: "{app}\RMBG.exe"; Tasks: "desktopicon"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"
Name: "startmenuicon"; Description: "{cm:CreateStartMenuIcon}"; GroupDescription: "{cm:AdditionalIcons}"
Name: "addassoc"; Description: "{cm:AddContextOption}"; GroupDescription: "{cm:AdditionalOptions}"; Flags: unchecked

[Files]
Source: "build\bin\RMBG.exe"; DestDir: "{app}"; Flags: ignoreversion

[Registry]
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbg"; ValueType: string; ValueName: ""; ValueData: "{cm:ContextQuick}"; Flags: uninsdeletekey; Check: WizardIsTaskSelected('addassoc')
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbg\command"; ValueType: string; ValueName: ""; ValueData: """{app}\RMBG.exe"" ""%1"""; Flags: uninsdeletekey; Check: WizardIsTaskSelected('addassoc')
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbg"; ValueType: string; ValueName: "ExtendedSubCommandsKey"; ValueData: ""; Flags: uninsdeletekey; Check: WizardIsTaskSelected('addassoc')
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbg"; ValueType: string; ValueName: "Position"; ValueData: "Bottom"; Flags: uninsdeletekey; Check: WizardIsTaskSelected('addassoc')

Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbgm"; ValueType: string; ValueName: ""; ValueData: "{cm:ContextPrecise}"; Flags: uninsdeletekey; Check: WizardIsTaskSelected('addassoc')
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbgm\command"; ValueType: string; ValueName: ""; ValueData: """{app}\RMBG.exe"" ""%1"" ""birefnet-massive"""; Flags: uninsdeletekey; Check: WizardIsTaskSelected('addassoc')
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbgm"; ValueType: string; ValueName: "ExtendedSubCommandsKey"; ValueData: ""; Flags: uninsdeletekey; Check: WizardIsTaskSelected('addassoc')
Root: HKCR; Subkey: "SystemFileAssociations\image\shell\rmbgm"; ValueType: string; ValueName: "Position"; ValueData: "Bottom"; Flags: uninsdeletekey; Check: WizardIsTaskSelected('addassoc')