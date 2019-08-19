unit ufmeUploadFile;

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
  TfmeUploadFile = class(TFrame)
    pb: TProgressBar;
    pnl: TPanel;
    lblFile: TLabel;
    lblInfo: TLabel;
  private
    var
      FFileService: TFileService;
      FError: string;
      FHost, FFileName, FWorkDir: string;
      FTotalSize, FReadSize: Int64;
    function GetStatusCode: Integer;
  public
    constructor Create(AOwner: TComponent; const AHost, AFileName, AWorkDir: string);
    destructor Destroy; override;
    property StatusCode: Integer read GetStatusCode;
    property Error: string read FError;
    function Upload(): Boolean;
    { Public declarations }
  end;

implementation

{$R *.fmx}
uses
  uTools;
{ TfmeUploadFile }

constructor TfmeUploadFile.Create(AOwner: TComponent; const AHost, AFileName, AWorkDir: string);
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
end;

destructor TfmeUploadFile.Destroy;
begin
  inherited;
  if Assigned(FFileService) then
    FreeAndNil(FFileService);
end;

function TfmeUploadFile.GetStatusCode: Integer;
begin
  Result := FFileService.StatusCode;
end;

function TfmeUploadFile.Upload: Boolean;
var
  FStream: TFileStream;
  LocalFileName, UploadUrl, Info: string;
  sSize, size, StartSize, Speed: Int64;
  StartTime, LastTime, Time: Cardinal;
begin
  Result := False;

  LocalFileName := TPath.Combine(FWorkDir, FFileName);

  if not FileExists(LocalFileName) then
    FError := 'file not found'
  else
  begin
    sSize := 0;
    try
      if FFileService.InfoHead(FFileName, sSize) or (FFileService.StatusCode = 404) then
      begin
        try
          FStream := TFileStream.Create(LocalFileName, fmOpenRead or fmShareExclusive);
          FTotalSize := FStream.Size;
          StartSize := sSize;

					// 严格比较可以对比两个文件的md5
          if (FTotalSize - sSize) = 0 then
          begin
            FError := 'server file has been upload';
          end
          else
          begin
            TThread.Synchronize(nil,
              procedure
              begin
                pb.Max := FTotalSize;
                pb.Value := StartSize;
                lblInfo.Text := '0 KB/s';
              end);

            StartTime := TThread.GetTickCount;

            UploadUrl := FFileService.UrlUpload;

            while (FTotalSize - sSize) > 0 do
            begin
              size := FTotalSize - sSize;
              if size = 0 then
              begin
                Break;
              end
              else if size < 0 then
              begin
                FError := 'server file size error';
                Exit;
              end
              else if size > FFileService.UploadChunkSize then
                size := FFileService.UploadChunkSize;

              if not FFileService.UploadChunk(FStream, UploadUrl, sSize, sSize + size - 1) then
              begin
                FError := FFileService.Error;
                Exit;
              end;

              Inc(sSize, size);

              Info := lblInfo.Text;
              Time := TThread.GetTickCount - StartTime; //计算用时
              if (Time - LastTime) > 800 then
              begin
                Speed := ((sSize - StartSize) * 1000) div Time; //计算速度
                Info := BytesToStr(Speed);

                LastTime := Time;
              end;

              TThread.Queue(nil,
                procedure
                begin
                  pb.Value := sSize;
                  lblInfo.Text := Info;
                end);
            end;

            if (FTotalSize - sSize) = 0 then
            begin
              FError := Format('upload file[%s] success', [FFileName]);
              Result := True;
            end;

            TThread.Queue(nil,
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
  end;
end;

end.

