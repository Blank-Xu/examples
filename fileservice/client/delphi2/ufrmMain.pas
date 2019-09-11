unit ufrmMain;

interface

uses
  System.SysUtils,
  System.Types,
  System.UITypes,
  System.Classes,
  System.Variants,
  System.IOUtils,
  FMX.Types,
  FMX.Controls,
  FMX.Forms,
  FMX.Graphics,
  FMX.Dialogs,
  FMX.StdCtrls,
  FMX.ScrollBox,
  FMX.Memo,
  FMX.Edit,
  FMX.Controls.Presentation,
  uFileService,
  ufmeDownloadFile,
  ufmeUploadFile;

type
  TfrmMain = class(TForm)
    pnl1: TPanel;
    edtHost: TEdit;
    edtFileName: TEdit;
    btnDownload: TButton;
    btnUpload: TButton;
    btnDelete: TButton;
    mmo: TMemo;
    pnlPB: TPanel;
    btnDownload2: TButton;
    pnlPB2: TPanel;
    pb: TProgressBar;
    pnl2: TPanel;
    lblInfo: TLabel;
    lblFile: TLabel;
    edtUsername: TEdit;
    edtPassword: TEdit;
    btnLogin: TButton;
    btnInfo: TButton;
    procedure btnDownloadClick(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure btnUploadClick(Sender: TObject);
    procedure btnDeleteClick(Sender: TObject);
    procedure btnDownload2Click(Sender: TObject);
    procedure btnLoginClick(Sender: TObject);
    procedure btnInfoClick(Sender: TObject);
  private
    const
      FWorkDir = 'files';
    var
      FToken: string;
      FTokenTime: TDateTime;
      FTotalSize: Int64;
    function CheckToken: Boolean;
    procedure ProgressCallback(Sender: TObject; Processed: Int64; SIZE: Int64; ContentLength: Int64; TimeStart: Cardinal);
  public
		{ Public declarations }
  end;

var
  frmMain: TfrmMain;

implementation
	{$R *.fmx}

uses
  uProgressFileStream,
  DateUtils;

procedure TfrmMain.btnDeleteClick(Sender: TObject);
var
  fme: TFileService;
  msg: string;
begin
  if not CheckToken then
    Exit;

  fme := TFileService.Create(nil, edtHost.Text, edtFileName.Text);
  try
    fme.Token := FToken;
    if fme.DeleteFile then
      msg := 'delete success'
    else
      msg := 'delete failed, err: ' + fme.Error;

    mmo.Lines.Add(msg);
  finally
    fme.Free;
  end;
end;

procedure TfrmMain.btnDownload2Click(Sender: TObject);
var
  fStream: TProgressFileStream;
  url, Msg: string;
  size: Int64;
  Host, FileName, LocalFileName: string;
  fme: TFileService;
begin
  if not CheckToken then
    Exit;

  Host := edtHost.Text;
  FileName := edtFileName.Text;
  LocalFileName := TPath.Combine(FWorkDir, FileName);

  pb.Value := 0;
  pb.Max := 1000;
  lblInfo.Text := '0 KB/s';
  lblFile.Text := edtFileName.Text;

  TThread.CreateAnonymousThread(
    procedure
    begin
      try
        try
          fme := TFileService.Create(nil, Host, FileName);
          fme.Token := FToken;
          if fme.InfoHead(FileName, FTotalSize) or (fme.StatusCode = 404) then
          begin
            try
              if FileExists(LocalFileName) then
                fStream := TProgressFileStream.Create(LocalFileName, fmOpenWrite or fmShareExclusive)
              else
                fStream := TProgressFileStream.Create(LocalFileName, fmCreate or fmShareExclusive);
              fStream.OnProgress := ProgressCallback;
							// 严格比较可以对比两个文件的md5
              size := FTotalSize - fStream.Size;
              if size = 0 then
              begin
                Msg := 'file have been download';
              end
              else
              begin
                TThread.Synchronize(nil,
                  procedure
                  begin
                    pb.Max := FTotalSize;
                    pb.Value := size;
                  end);

                url := fme.UrlDownload;
                while (FTotalSize - fStream.Size) > 0 do
                begin
                  size := FTotalSize - fStream.Size;
                  if size = 0 then
                  begin
                    Break;
                  end
                  else if size < 0 then
                  begin
                    Msg := 'server file size error';
                    Break;
                  end
                  else if size > fme.DownloadChunkSize then
                    size := fme.DownloadChunkSize;

                  if not fme.DownloadChunk(fStream, url, fStream.Size, fStream.Size + size - 1) then
                  begin
                    Msg := fme.Error;
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
            Msg := fme.Error;
        except
          on E: Exception do
            Msg := E.Message;
        end;
      finally
        if Assigned(fme) then
          FreeAndNil(fme);
      end;

      TThread.Queue(nil,
        procedure
        begin
          lblInfo.Text := '0 KB/s';
          pb.Value := pb.Max;
          mmo.Lines.Add(Msg);
        end);
    end).Start;
end;

procedure TfrmMain.btnDownloadClick(Sender: TObject);
begin
  if not CheckToken then
    Exit;

  TThread.CreateAnonymousThread(
    procedure
    var
      fmeDownload: TfmeDownloadFile;
      Msg: string;
    begin
      fmeDownload := TfmeDownloadFile.Create(pnlPB, edtHost.Text, FToken, edtFileName.Text, FWorkDir);
      fmeDownload.Parent := pnlPB;
      fmeDownload.Visible := True;

      TThread.Synchronize(nil,
        procedure
        begin
          pnlPB.Repaint;
        end);

      try
        if not fmeDownload.Download then
          Msg := 'download failed, err: ' + fmeDownload.Error
        else
          Msg := 'download success';

        TThread.Synchronize(nil,
          procedure
          begin
            mmo.Lines.Add(Msg);
            fmeDownload.Visible := False;
          end);

      finally
        Sleep(1000);

        if Assigned(fmeDownload) then
          FreeAndNil(fmeDownload);
      end;
    end).Start;
end;

procedure TfrmMain.btnInfoClick(Sender: TObject);
var
  fs: TFileService;
  MStream: TMemoryStream;
  Stream: TStringStream;
begin
  if not CheckToken then
    Exit;

  fs := TFileService.Create(nil, edtHost.Text);
  MStream := TMemoryStream.Create;
  Stream := TStringStream.Create;
  try
    fs.Token := FToken;
    if fs.Info(edtFileName.Text, MStream, True) then
    begin
      MStream.Position := 0;
      Stream.LoadFromStream(MStream);
      mmo.Lines.Add(Stream.DataString);
    end
    else
      mmo.Lines.Add(Format('code: %d, msg: %s', [fs.StatusCode, fs.Error]));
  finally
    if Assigned(MStream) then
      FreeAndNil(MStream);
    if Assigned(Stream) then
      FreeAndNil(Stream);
    if Assigned(fs) then
      FreeAndNil(fs);
  end;
end;

procedure TfrmMain.btnLoginClick(Sender: TObject);
var
  fs: TFileService;
  Username, Password: string;
begin
  if (Length(edtUsername.Text) = 0) or (Length(edtPassword.Text) = 0) then
  begin
    ShowMessage('请输入完整的用户名和密码！');
    Exit;
  end;

  fs := TFileService.Create(nil, edtHost.Text);
  try
    if fs.Login(edtUsername.Text, edtPassword.Text) then
    begin
      FToken := fs.Token;
      FTokenTime := Now;
      mmo.Lines.Add('login success, token: ' + FToken);
    end
    else
    begin
      mmo.Lines.Add('login failed, err: ' + fs.Error);
    end;
  finally
    if Assigned(fs) then
      FreeAndNil(fs);
  end;
end;

procedure TfrmMain.btnUploadClick(Sender: TObject);
begin
  if not CheckToken then
    Exit;

  TThread.CreateAnonymousThread(
    procedure
    var
      fmeUpload: TfmeUploadFile;
      Msg: string;
    begin
      fmeUpload := TfmeUploadFile.Create(pnlPB, edtHost.Text, FToken, edtFileName.Text, FWorkDir);
      fmeUpload.Parent := pnlPB;
      fmeUpload.Visible := True;

      TThread.Synchronize(nil,
        procedure
        begin
          pnlPB.Repaint;
        end);

      try
        if not fmeUpload.Upload then
          Msg := 'upload failed, err: ' + fmeUpload.Error
        else
          Msg := 'upload success';

        TThread.Synchronize(nil,
          procedure
          begin
            mmo.Lines.Add(Msg);
            fmeUpload.Visible := False;
          end);

      finally
        Sleep(1000);

        if Assigned(fmeUpload) then
          FreeAndNil(fmeUpload);
      end;
    end).Start;
end;

function TfrmMain.CheckToken: Boolean;
begin
  Result := True;
	// token少于5分钟重新登录获取，避免下载大文件超时
  if (Length(FToken) = 0) or (SecondsBetween(Now, FTokenTime) >= 25 * 30) then
    Result := False;
  mmo.Lines.Add('token expire seconds: ' + IntToStr(25 * 30 - SecondsBetween(Now, FTokenTime)));
end;

procedure TfrmMain.FormCreate(Sender: TObject);
begin
  if not DirectoryExists(FWorkDir) then
    MkDir(FWorkDir);
end;

procedure TfrmMain.ProgressCallback(Sender: TObject; Processed, Size, ContentLength: Int64; TimeStart: Cardinal);
begin
  TThread.Queue(nil,
    procedure
    begin
      if Processed > pb.Max then
        pb.Value := pb.Max
      else
        pb.Value := Processed;
    end);
end;

end.

