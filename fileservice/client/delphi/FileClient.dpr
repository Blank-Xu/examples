program FileClient;

uses
  Vcl.Forms,
  uClient in 'uClient.pas' {frmClient},
  uFile in 'uFile.pas';

{$R *.res}

begin
  Application.Initialize;
  Application.MainFormOnTaskbar := True;
  Application.CreateForm(TfrmClient, frmClient);
  Application.Run;
end.
