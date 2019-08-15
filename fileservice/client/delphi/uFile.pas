unit uFile;

interface

uses
	System.Classes, System.SysUtils, System.Net.URLClient, System.Net.HttpClient,
  System.Net.HttpClientComponent;

type
  TFileService = class
  private
    const
      bytes = 'bytes';
      range = 'bytes=%d-%d';
      requestTimeout = 20;
      infoUrl = '%s/info?filename=%s&checkMd5=%s';
      uploadUrl = '%s/upload?filename=%s';
      downloadUrl = '%s/download?filename=%s';
    var
      httpClient: TNetHTTPClient;
      Host: string;
      FileName: string;
    function OpenFile(fileName: string): Boolean;
  public
    constructor Create; override;
    destructor Destroy;
    function FileInfo(fileName: string; checkMd5: Boolean): Boolean;
    function DownloadFile(fileName: string): Boolean;
    function UploadFile(fileName: string): Boolean;
  end;

implementation

{ TFileService }

constructor TFileService.Create;
begin
  inherited;
  httpClient := TNetHTTPClient.Create;
  httpClient.ConnectionTimeout := 3;
  httpClient.ResponseTimeout := requestTimeout;
  httpClient.CustomHeaders['Keep-Alive'] := '60';
end;

destructor TFileService.Destroy;
begin
  if Assigned(httpClient) then
    httpClient.Free;
end;

function TFileService.DownloadFile(fileName: string): Boolean;
begin

  httpClient.CustomHeaders['Range'] := range;
end;

function TFileService.FileInfo(fileName: string; checkMd5: Boolean): Boolean;
var
  url: string;
	Stream: TMemoryStream;
	AHeaders: TNetHeaders  ;
	Resp :IHTTPResponse;
begin
	url := Format(infoUrl, [Host, fileName, checkMd5]);
	Stream := TMemoryStream.Create;

	Resp := httpClient.Get(url,Stream,nil);
	if Resp.StatusCode <> 200  then
	begin
		Result = False;
	end;
end;

function TFileService.OpenFile(fileName: string): Boolean;
begin

end;

function TFileService.UploadFile(fileName: string): Boolean;
begin

end;

end.

