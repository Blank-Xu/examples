unit uClient;

interface

uses
  Winapi.Windows,
  Winapi.Messages,
  System.SysUtils,
  System.Variants,
  System.Classes,
  Vcl.Graphics,
  Vcl.Controls,
  Vcl.Forms,
  Vcl.Dialogs,
  uFileService,
  Vcl.StdCtrls,
  Vcl.ExtCtrls;

type
  TfrmClient = class(TForm)
    mmo: TMemo;
    pnl1: TPanel;
    btn1: TButton;
    edtHost: TEdit;
    btn2: TButton;
    btn3: TButton;
    btn4: TButton;
    edtFilename: TEdit;
    procedure btn1Click(Sender: TObject);
    procedure btn2Click(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure btn3Click(Sender: TObject);
    procedure FormDestroy(Sender: TObject);
    procedure btn4Click(Sender: TObject);
  private
    var
      FFileService: TFileService;
    procedure addLine;
  public
    { Public declarations }
  end;

var
  frmClient: TfrmClient;

implementation

{$R *.dfm}

procedure TfrmClient.addLine;
begin
  if mmo.Lines.Count > 500 then
    mmo.Lines.Clear;

  mmo.Lines.Add('');
  mmo.Lines.Add('----------------------------------');
end;

procedure TfrmClient.btn1Click(Sender: TObject);
var
  size: Int64;
begin
  addLine;

  FFileService.Host := edtHost.Text;

  if FFileService.InfoHead(edtFilename.Text, size) then
    mmo.Lines.Add(Format('size: %d, mod_time: %s', [size, FFileService.Response.LastModified]))
  else
    mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));
end;

procedure TfrmClient.btn2Click(Sender: TObject);
var
  stream: TStringStream;
begin
  addLine;

  FFileService.Host := edtHost.Text;

  stream := TStringStream.Create;
  try
    if FFileService.Info(edtFilename.Text, stream, True) then
    begin
      stream.Position := 0;
      mmo.Lines.Add(stream.DataString);
    end
    else
      mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));
  finally
    if Assigned(stream) then
      FreeAndNil(stream);
  end;
end;

procedure TfrmClient.btn3Click(Sender: TObject);
begin
  addLine;

  FFileService.Host := edtHost.Text;

  if FFileService.DownloadFile(edtFilename.Text) then
  begin
    mmo.Lines.Add(Format('download file[%s] sucess', [edtFilename.Text]));
  end
  else
    mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));
end;

procedure TfrmClient.btn4Click(Sender: TObject);
begin
  addLine;

  FFileService.Host := edtHost.Text;

  if FFileService.UploadFile(edtFilename.Text) then
  begin
    mmo.Lines.Add(Format('upload file[%s] sucess', [edtFilename.Text]));
  end
  else
    mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));
end;

procedure TfrmClient.FormCreate(Sender: TObject);
begin
  ReportMemoryLeaksOnShutdown := True;
  FFileService := TFileService.Create('');
  FFileService.WorkDir := 'files';
end;

procedure TfrmClient.FormDestroy(Sender: TObject);
begin
  if Assigned(FFileService) then
    FreeAndNil(FFileService);
end;

end.

