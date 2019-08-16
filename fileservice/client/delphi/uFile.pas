unit uFile;

interface

uses
  System.Classes,
  System.SysUtils,
  System.Net.URLClient,
  System.Net.HttpClient,
  System.Net.HttpClientComponent;

type
  TFileService = class
  private
    const
      CONNECTION_TIMEOUT = 30;
      RESPONSE_TIMEOUT = 30;
      RANGE_BYTES = 'bytes=%d-%d';
      INFO_URL = '%s/info?filename=%s&md5=%s';
      UPLOAD_URL = '%s/upload?filename=%s';
      DOWNLOAD_URL = '%s/download?filename=%s';
    var
      FHttpClient: TNetHTTPClient;
      FHost: string;
      FFileName: string;
    function OpenFile(const fileName: string; const mode: Word; fStream: TFileStream; err: string): Boolean;
  public
    constructor Create(host: string);
    destructor Destroy;
    function FileInfo(fileName: string; md5: Boolean = False): Boolean;
    function DownloadFile(fileName: string): Boolean;
    function UploadFile(fileName: string): Boolean;
  end;

implementation

{ TFileService }

constructor TFileService.Create(host: string);
begin
  FHost := host;

  FHttpClient := TNetHTTPClient.Create(nil);
  FHttpClient.ConnectionTimeout := CONNECTION_TIMEOUT;
  FHttpClient.ResponseTimeout := RESPONSE_TIMEOUT;
  FHttpClient.CustomHeaders['Keep-Alive'] := '60';
end;

destructor TFileService.Destroy;
begin
  if Assigned(FHttpClient) then
    FHttpClient.Free;
end;

function TFileService.DownloadFile(fileName: string): Boolean;
begin

//	FHttpClient.CustomHeaders['Range'] := range;
end;

function TFileService.FileInfo(fileName: string; md5: Boolean = False): Boolean;
var
  url: string;
  Stream: TMemoryStream;
  Resp: IHTTPResponse;
begin
  Result := False;
  url := Format(INFO_URL, [FHost, fileName, BoolToStr(md5, True)]);

  Stream := TMemoryStream.Create;
  try
    try
      Resp := FHttpClient.Get(url, Stream, nil);
      if Resp.StatusCode <> 200 then
      begin

      end;

      Result := True;
    except
      on E: Exception do


    end;
  finally
    Stream.Free;
  end;
end;

function TFileService.OpenFile(const fileName: string; const mode: Word; fStream: TFileStream; err: string): Boolean;
begin
  Result := False;
  if not Assigned(fStream) then
  begin
    try
      fStream := TFileStream.Create(fileName, mode);
      Result := True;
    except
      on E: Exception do
        err := E.Message;
    end;
  end;
end;

function TFileService.UploadFile(fileName: string): Boolean;
begin

end;

end.

