unit uClient;

interface

uses
  Winapi.Windows,
  Winapi.Messages,
  System.SysUtils,
  System.Variants,
  System.Classes,
  System.IOUtils,
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
    btnInfoHead: TButton;
    edtHost: TEdit;
    btnInfo: TButton;
    btnDownload: TButton;
    btnUpload: TButton;
    edtFilename: TEdit;
    pb: TProgressBar;
    procedure btnInfoHeadClick(Sender: TObject);
    procedure btnInfoClick(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure btnDownloadClick(Sender: TObject);
    procedure FormDestroy(Sender: TObject);
    procedure btnUploadClick(Sender: TObject);
    procedure FormCloseQuery(Sender: TObject; var CanClose: Boolean);
  private
    const
      FWorkDir = 'files';
    var
      FCanClose: Boolean;
      FFileService: TFileService;
      FTotalSize: Int64;
      FReadSize: Int64;
      FNeedProcess: Boolean;
      FIsDownload: Boolean;
    procedure StartProcessing;
    procedure EndProcessing;
    procedure AddLine;
    procedure SetTotalSize(const size: Int64; const Download: Boolean = True);
    procedure DownloadData(const Sender: TObject; AContentLength, AReadCount: Int64; var Abort: Boolean);
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

procedure TfrmClient.btnInfoHeadClick(Sender: TObject);
begin
  StartProcessing;

  if FFileService.InfoHead(edtFilename.Text, FTotalSize) then
    mmo.Lines.Add(Format('size: %d, mod_time: %s', [FTotalSize, FFileService.Response.LastModified]))
  else
    mmo.Lines.Add(Format('code: %d, msg: %s', [FFileService.StatusCode, FFileService.Error]));

  EndProcessing;
end;

procedure TfrmClient.btnInfoClick(Sender: TObject);
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

procedure TfrmClient.btnDownloadClick(Sender: TObject);
var
  fStream: TFileStream;
  url, Msg: string;
  sSize: Int64;
  size: Int64;
  FileName, LocalFileName: string;
begin
  StartProcessing;
  FileName := edtFilename.Text;
  LocalFileName := TPath.Combine(FWorkDir, FileName);
  FFileService.FileName := FileName;
  sSize := 0;
  try
    if FFileService.InfoHead(FileName, sSize) or (FFileService.StatusCode = 404) then
    begin
      SetTotalSize(sSize);
      try
        if FileExists(LocalFileName) then
          fStream := TFileStream.Create(LocalFileName, fmOpenWrite or fmShareExclusive)
        else
          fStream := TFileStream.Create(LocalFileName, fmCreate or fmShareExclusive);
					// 严格比较可以对比两个文件的md5
        if (sSize - fStream.Size) = 0 then
        begin
          Msg := 'file have been download';
        end
        else
        begin
          url := FFileService.UrlDownload;
          while (sSize - fStream.Size) > 0 do
          begin
            size := sSize - fStream.Size;
            if size = 0 then
            begin
              Break;
            end
            else if size < 0 then
            begin
              Msg := 'server file size error';
              Break;
            end
            else if size > FFileService.DownloadChunkSize then
              size := FFileService.DownloadChunkSize;

            if not FFileService.DownloadChunk(fStream, url, fStream.Size, fStream.Size + size - 1) then
            begin
              Msg := FFileService.Error;
              Break;
            end;
          end;
          Msg := Format('download file[%s] success', [FileName]);
        end;
      finally
        if Assigned(fStream) then
          FreeAndNil(fStream);
      end;
    end
    else
      Msg := FFileService.Error;
  except
    on E: Exception do
      Msg := E.Message;
  end;
  mmo.Lines.Add(Msg);

  EndProcessing;
end;

procedure TfrmClient.btnUploadClick(Sender: TObject);
var
  fStream: TFileStream;
  url, Msg: string;
  sSize: Int64;
  size: Int64;
  FileName, LocalFileName: string;
begin
  StartProcessing;
  FileName := edtFilename.Text;
  LocalFileName := TPath.Combine(FWorkDir, FileName);
  FFileService.FileName := FileName;
  if not FileExists(LocalFileName) then
    Msg := 'file not found'
  else
  begin
    sSize := 0;
    try
      if FFileService.InfoHead(FileName, sSize) or (FFileService.StatusCode = 404) then
      begin
        try
          fStream := TFileStream.Create(LocalFileName, fmOpenRead or fmShareExclusive);
          FTotalSize := fStream.Size;
          SetTotalSize(FTotalSize);
					// 严格比较可以对比两个文件的md5
          if (FTotalSize - sSize) = 0 then
          begin
            Msg := 'server file has been upload';
          end
          else
          begin
            url := FFileService.UrlUpload;
            while (FTotalSize - sSize) > 0 do
            begin
              size := FTotalSize - sSize;
              if size = 0 then
              begin
                Break;
              end
              else if size < 0 then
              begin
                Msg := 'server file size error';
                Break;
              end
              else if size > FFileService.UploadChunkSize then
                size := FFileService.UploadChunkSize;

              if not FFileService.UploadChunk(fStream, url, sSize, sSize + size - 1) then
              begin
                Msg := FFileService.Error;
                Break;
              end;
              Inc(sSize, size);

              if FNeedProcess then
              begin
                pb.Position := sSize;
                pb.Refresh;
              end;
            end;
            Msg := Format('upload file[%s] success', [FileName]);
          end;
        finally
          if Assigned(fStream) then
            FreeAndNil(fStream);
        end;
      end
      else
        Msg := FFileService.Error;
    except
      on E: Exception do
        Msg := E.Message;
    end;
  end;
  mmo.Lines.Add(Msg);

  EndProcessing;
end;

procedure TfrmClient.EndProcessing;
begin
  FTotalSize := 0;
  FNeedProcess := False;
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

  FFileService := TFileService.Create(nil);
  FFileService.WorkDir := FWorkDir;
  FFileService.OnReceiveData := DownloadData;
end;

procedure TfrmClient.FormDestroy(Sender: TObject);
begin
  if Assigned(FFileService) then
    FreeAndNil(FFileService);
end;

procedure TfrmClient.DownloadData(const Sender: TObject; AContentLength, AReadCount: Int64; var Abort: Boolean);
begin
  if FIsDownload and FNeedProcess then
  begin
    if AContentLength = AReadCount then
      Inc(FReadSize, AReadCount);
    pb.Position := FReadSize + AReadCount;
    pb.Refresh;
  end;
end;

procedure TfrmClient.SetTotalSize(const Size: Int64; const Download: Boolean = True);
begin
  FTotalSize := Size;
  if FTotalSize > 10 * 1024 * 1024 then
  begin
    pb.Max := FTotalSize;
    FNeedProcess := True;
  end;
  if Download then
    FIsDownload := True;
end;

procedure TfrmClient.StartProcessing;
begin
  FCanClose := False;

  FFileService.Host := edtHost.Text;
  pb.Position := 0;

  FReadSize := 0;

  AddLine;
end;

end.

