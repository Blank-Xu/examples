program FileService;

uses
  System.StartUpCopy,
  FMX.Forms,
  ufrmMain in 'ufrmMain.pas' {frmMain},
  uFileService in 'uFileService.pas',
  ufmeDownloadFile in 'ufmeDownloadFile.pas' {fmeDownloadFile: TFrame},
	ufmeUploadFile in 'ufmeUploadFile.pas' {fmeUploadFile: TFrame},
	uTools in 'uTools.pas';

{$R *.res}

begin
{IFDEF DEBUG}
	ReportMemoryLeaksOnShutdown := True;
{ENDIF}
  Application.Initialize;
  Application.CreateForm(TfrmMain, frmMain);
  Application.Run;
end.

