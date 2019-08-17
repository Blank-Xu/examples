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
  Vcl.StdCtrls,
  Vcl.ExtCtrls,
  Vcl.ComCtrls,
  uFileService;

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
    pb: TProgressBar;
    procedure btn1Click(Sender: TObject);
    procedure btn2Click(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure btn3Click(Sender: TObject);
    procedure FormDestroy(Sender: TObject);
    procedure btn4Click(Sender: TObject);
    procedure FormCloseQuery(Sender: TObject; var CanClose: Boolean);
  private
    var
      FCanClose: Boolean;
      FFileService: TFileService;
      FTotalSize: Int64;
      FReadSize: Int64;
      FReadPosition: Boolean;
    procedure StartProcessing;
    procedure EndProcessing;
    procedure AddLine;
    procedure ReceiveData(const Sender: TObject; AContentLength, AReadCount: Int64; var Abort: Boolean);
  public
    { Public declarations }
  end;

var
  frmClient: TfrmClient;

implementation

{$R *.dfm}

procedure TfrmClient.AddLine;
begin
  if mmo.Lines.Count > 500 then
    mmo.Lines.Clear;

  mmo.Lines.Add('');
  mmo.Lines.Add('----------------------------------');
end;

procedure TfrmClient.btn1Click(Sender: TObject);
begin
  StartProcessing;

  if FFileService.InfoHead(edtFilename.Text, FTotalSize) then
    mmo.Lines.Add(Format('size: %d, mod_time: %s', [FTotalSize, FFileService.Response.LastModified]))
  else
    mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));

  EndProcessing;
end;

procedure TfrmClient.btn2Click(Sender: TObject);
var
  Stream: TStringStream;
begin
  StartProcessing;

  Stream := TStringStream.Create;
  try
    if FFileService.Info(edtFilename.Text, Stream, True) then
    begin
      Stream.Position := 0;
      mmo.Lines.Add(Stream.DataString);
    end
    else
      mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));
  finally
    if Assigned(Stream) then
      FreeAndNil(Stream);
  end;

  EndProcessing;
end;

procedure TfrmClient.btn3Click(Sender: TObject);
var
  msg: string;
begin
  StartProcessing;

  if FFileService.DownloadFile(edtFilename.Text) then
  begin
    mmo.Lines.Add(Format('download file[%s] sucess', [edtFilename.Text]));
  end
  else
    mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));

  EndProcessing;
end;

procedure TfrmClient.btn4Click(Sender: TObject);
var
	fs: TFileService;
	msg: string;
begin
	StartProcessing;

	TThread.CreateAnonymousThread(
		procedure
		var
			fs: TFileService;
		begin
      if FFileService.UploadFile(edtFilename.Text) then
        msg := Format('upload file[%s] sucess', [edtFilename.Text])
      else
        msg := Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]);

      TThread.Synchronize(nil,
        procedure
        begin
					mmo.Lines.Add(msg)
        end)
		end).Start;
//	if FFileService.UploadFile(edtFilename.Text) then
//	begin
//		mmo.Lines.Add(Format('upload file[%s] sucess', [edtFilename.Text]));
//	end
//	else
//		mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));

  EndProcessing;
end;

procedure TfrmClient.EndProcessing;
begin
  FReadPosition := False;
  FCanClose := True;
end;

procedure TfrmClient.FormCloseQuery(Sender: TObject; var CanClose: Boolean);
begin
  CanClose := FCanClose;
end;

procedure TfrmClient.FormCreate(Sender: TObject);
begin
  ReportMemoryLeaksOnShutdown := True;
  FCanClose := True;

  FFileService := TFileService.Create('');
  FFileService.WorkDir := 'files';
  FFileService.OnReceiveData := ReceiveData;
end;

procedure TfrmClient.FormDestroy(Sender: TObject);
begin
  if Assigned(FFileService) then
    FreeAndNil(FFileService);
end;

procedure TfrmClient.ReceiveData(const Sender: TObject; AContentLength, AReadCount: Int64; var Abort: Boolean);
begin
  if FReadPosition then
  begin
    if AContentLength = AReadCount then
      Inc(FReadSize, AReadCount);
    pb.Position := FReadSize + AReadCount;
    pb.Refresh;
  end;
end;

procedure TfrmClient.StartProcessing;
begin
  FCanClose := False;

  FFileService.Host := edtHost.Text;
  pb.Position := 0;

  FReadSize := 0;

  if FTotalSize > 10 * 1024 * 1024 then
  begin
    pb.Max := FTotalSize;
    FReadPosition := True;
  end;

  AddLine;
end;

end.

