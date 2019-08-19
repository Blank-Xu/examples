unit ufmeDownloadFile;

interface

uses
  System.SysUtils,
  System.Types,
  System.UITypes,
  System.Classes,
  System.Variants,
  System.IOUtils,
  FMX.Types,
  FMX.Graphics,
  FMX.Controls,
  FMX.Forms,
  FMX.Dialogs,
  FMX.StdCtrls,
  FMX.Controls.Presentation,
  uFileService;

type
  TfmeDownloadFile = class(TFrame)
    pb: TProgressBar;
    lblFile: TLabel;
    lblInfo: TLabel;
    pnl: TPanel;
  private
    var
      FFileService: TFileService;
      FError: string;
      FHost, FFileName, FWorkDir: string;
      FTotalSize, FReadSize, FStartSize: Int64;
      FStartTime, FLastTime: Cardinal;
    function GetStatusCode: Integer;
    procedure SetStatus(ASize: Int64);
    procedure OnReceiveData(const Sender: TObject; AContentLength, AReadCount: Int64; var Abort: Boolean);
  public
    constructor Create(AOwner: TComponent; const AHost, AFileName, AWorkDir: string);
    destructor Destroy; override;
    property StatusCode: Integer read GetStatusCode;
    property Error: string read FError;
    function Download(): Boolean;
  end;

implementation

{$R *.fmx}

uses
  uTools;

{ TfmeDownloadFile }

constructor TfmeDownloadFile.Create(AOwner: TComponent; const AHost, AFileName, AWorkDir: string);
begin
  inherited Create(AOwner);

  FHost := AHost;
  FFileName := AFileName;
  FWorkDir := AWorkDir;

  TThread.Queue(nil,
    procedure
    begin
      pb.Max := 0;
      pb.Value := 0;
      lblFile.Text := AFileName;
    end);

  FFileService := TFileService.Create(AOwner, AHost, AFileName);
  FFileService.OnReceiveData := OnReceiveData;
end;

destructor TfmeDownloadFile.Destroy;
begin
  inherited;
  if Assigned(FFileService) then
    FreeAndNil(FFileService);
end;

function TfmeDownloadFile.Download: Boolean;
var
  FStream: TFileStream;
  Size: Int64;
  DownloadUrl: string;
  LocalFileName: string;
begin
  Result := False;

  FTotalSize := 0;
  try
    if FFileService.InfoHead(FTotalSize) or (FFileService.StatusCode = 404) then
    begin
      SetStatus(FTotalSize);
      LocalFileName := TPath.Combine(FWorkDir, FFileName);
      try
        if FileExists(LocalFileName) then
          FStream := TFileStream.Create(LocalFileName, fmOpenWrite or fmShareExclusive)
        else
          FStream := TFileStream.Create(LocalFileName, fmCreate or fmShareExclusive);

        FStartSize := FTotalSize - FStream.Size;
				// 严格比较可以对比两个文件的md5
        if FStartSize = 0 then
        begin
          FError := 'file have been download';
        end
        else
        begin
          TThread.Synchronize(nil,
            procedure
            begin
              pb.Max := FTotalSize;
              pb.Value := FStartSize;
              lblInfo.Text := '0 KB/s';
            end);

          DownloadUrl := FFileService.UrlDownload;

          while (FTotalSize - FStream.Size) > 0 do
          begin
            Size := FTotalSize - FStream.Size;
            if Size = 0 then
            begin
              Break;
            end
            else if Size < 0 then
            begin
              FError := 'server file size error';
              Break;
            end
            else if Size > FFileService.DownloadChunkSize then
              Size := FFileService.DownloadChunkSize;

            if not FFileService.DownloadChunk(FStream, DownloadUrl, FStream.Size, FStream.Size + Size - 1) then
            begin
              FError := FFileService.Error;
              Break;
            end;
          end;

          if (FTotalSize - FStream.Size) = 0 then
          begin
            FError := Format('download file[%s] success', [FFileName]);
            Result := True;
          end;

          TThread.Synchronize(nil,
            procedure
            begin
              lblInfo.Text := '0 KB/s';
            end);
        end;
      finally
        if Assigned(FStream) then
          FreeAndNil(FStream);
      end;
    end
    else
      FError := FFileService.Error;
  except
    on E: Exception do
      FError := E.Message;
  end;

  TThread.Queue(nil,
    procedure
    begin
      lblInfo.Text := FError;
    end);
end;

function TfmeDownloadFile.GetStatusCode: Integer;
begin
  Result := FFileService.StatusCode;
end;

procedure TfmeDownloadFile.OnReceiveData(const Sender: TObject; AContentLength, AReadCount: Int64; var Abort: Boolean);
var
  Time: Cardinal;
  Value, Speed: Int64;
  Info: string;
begin
  Info := lblInfo.Text;

  if AContentLength = AReadCount then
  begin
    Inc(FReadSize, AReadCount);
    Value := FStartSize + FReadSize;
  end
  else
    Value := FStartSize + FReadSize + AReadCount;

  if Value >= FTotalSize then
  begin
    Info := '0 KB/s';
  end
  else
  begin
    Info := lblInfo.Text;
    Time := TThread.GetTickCount - FStartTime; //计算用时
    if (Time - FLastTime) > 900 then
    begin
      Speed := ((Value - FStartSize) * 1000) div Time; //计算速度
      Info := BytesToStr(Speed);

      FLastTime := Time;
    end;
  end;

  TThread.Queue(nil,
    procedure
    begin
      pb.Value := Value;
      lblInfo.Text := Info;
    end);
end;

procedure TfmeDownloadFile.SetStatus(ASize: Int64);
begin
  pb.Max := ASize;
  FStartTime := TThread.GetTickCount;
  FLastTime := 0;
end;

end.

