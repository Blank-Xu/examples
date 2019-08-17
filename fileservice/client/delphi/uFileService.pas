unit uFileService;

interface

uses
  System.Classes,
  System.SysUtils,
  System.IOUtils,
  System.Net.URLClient,
  System.Net.HttpClient,
  System.Net.HttpClientComponent;

type
  TFileService = class(TNetHTTPClient)
  private
    const
      CONNECTION_TIMEOUT = 30;
      RESPONSE_TIMEOUT = 30;
      RANGE_BYTES = 'bytes=%d-%d';
      URL_INFO_HEAD = '%s/info?filename=%s';
      URL_INFO = '%s/info?filename=%s';
      URL_INFO_MD5 = '%s/info?filename=%s&md5=true';
      URL_UPLOAD = '%s/upload?filename=%s';
      URL_DOWNLOAD = '%s/download?filename=%s';
      URL_DELETE = '%s/delete?filename=%s';
    var
      FUploadChunkSize: Int64;
      FDownloadChunkSize: Int64;
      FResponse: IHTTPResponse;
      FHost: string;
      FWorkDir: string;
      FFileName: string;
      FError: string;
    procedure SetWorkDir(const dir: string);
    function GetUrlUpload(): string;
    function GetUrlDownload(): string;
			// »ñÈ¡×´Ì¬
    function GetStatusCode: Integer;
  public
    constructor Create(AOwner: TComponent); overload;
    constructor Create(const cHost: string); overload;
    constructor Create(const cHost, cFileName: string); overload;
    destructor Destroy; override;
    property UploadChunkSize: Int64 read FUploadChunkSize write FUploadChunkSize;
    property DownloadChunkSize: Int64 read FDownloadChunkSize write FDownloadChunkSize;
    property Host: string read FHost write FHost;
    property WorkDir: string read FWorkDir write SetWorkDir;
    property FileName: string read FFileName write FFileName;
    property UrlUpload: string read GetUrlUpload;
    property UrlDownload: string read GetUrlDownload;
    property Response: IHTTPResponse read FResponse;
    property StatusCode: Integer read GetStatusCode;
    property Error: string read FError;
		// for get file size and mod time
    function InfoHead(const FileName: string; var Size: Int64): Boolean;
    function Info(const FileName: string; fileInfoStream: TMemoryStream; Md5: Boolean = False): Boolean;
    function DownloadFile(const FileName: string): Boolean;
    function UploadFile(const FileName: string): Boolean;
    function DeleteFile(const FileName: string): Boolean;
    function DownloadChunk(FileStream: TFileStream; const url: string; const Offset, Size: Int64): Boolean;
    function UploadChunk(FileStream: TFileStream; const url: string; const Offset, Size: Int64): Boolean;
  end;

implementation

{ TFileService }

constructor TFileService.Create(const cHost, cFileName: string);
begin
  Create(nil);

  FHost := cHost;
  FFileName := cFileName;
end;

constructor TFileService.Create(const cHost: string);
begin
  Create(nil);

  FHost := cHost;
end;

constructor TFileService.Create(AOwner: TComponent);
begin
  inherited Create(nil);

  FUploadChunkSize := 4 * 1024 * 1024;
  FDownloadChunkSize := 4 * 1024 * 1024;

  HandleRedirects := True;
  UserAgent := 'client 1.0';
  ConnectionTimeout := CONNECTION_TIMEOUT;
  ResponseTimeout := RESPONSE_TIMEOUT;
  CustomHeaders['Keep-Alive'] := '60';
end;

function TFileService.DeleteFile(const fileName: string): Boolean;
var
  RespStream: TMemoryStream;
  SStream: TStringStream;
  url: string;
begin
  Result := False;
  url := Format(URL_DELETE, [FHost, fileName]);
  RespStream := TMemoryStream.Create;
  try
    try
      FResponse := Self.Post(url, nil, RespStream, nil);
      if FResponse.StatusCode = 200 then
        Exit(True)
      else if Assigned(RespStream) then
      begin
        SStream := TStringStream.Create;
        try
          SStream.LoadFromStream(RespStream);
          FError := SStream.DataString;
        finally
          SStream.Free;
        end;
      end;
      FError := 'file not found';
    except
      on E: Exception do
        FError := E.Message;
    end;
  finally
    if Assigned(RespStream) then
      FreeAndNil(RespStream);
  end;
end;

destructor TFileService.Destroy;
begin
  inherited;
end;

function TFileService.DownloadChunk(fileStream: TFileStream; const url: string; const offset, size: Int64): Boolean;
var
  AHeaders: TNetHeaders;
  RespStream: TMemoryStream;
  SStream: TStringStream;
begin
  Result := False;

  SetLength(AHeaders, Length(AHeaders) + 1);
  AHeaders[High(AHeaders)] := TNameValuePair.Create('Range', Format(RANGE_BYTES, [offset, size]));

  RespStream := TMemoryStream.Create;
  try
    try
      fileStream.Position := offset;

      FResponse := Self.Get(url, RespStream, AHeaders);
      if (FResponse.StatusCode = 200) or (FResponse.StatusCode = 206) then
      begin
        fileStream.CopyFrom(RespStream, 0);
        Exit(True);
      end
      else if Assigned(RespStream) then
      begin
        SStream := TStringStream.Create;
        try
          SStream.LoadFromStream(RespStream);
          FError := SStream.DataString;
        finally
          if Assigned(SStream) then
            FreeAndNil(SStream);
        end;
      end;
    except
      on E: Exception do
        FError := E.Message
    end;
  finally
    if Assigned(RespStream) then
      FreeAndNil(RespStream);
  end;
end;

function TFileService.DownloadFile(const fileName: string): Boolean;
var
  FStream: TFileStream;
  totalSize, size: Int64;
  url, lFileName: string;
begin
  Result := False;

  totalSize := 0;
  size := 0;

  if InfoHead(fileName, totalSize) then
  begin
    url := Format(URL_DOWNLOAD, [FHost, fileName]);
    lFileName := TPath.Combine(FWorkDir, fileName);
    try
      try
        if FileExists(lFileName) then
          FStream := TFileStream.Create(lFileName, fmOpenWrite or fmShareExclusive)
        else
          FStream := TFileStream.Create(lFileName, fmCreate or fmShareExclusive);

        if (totalSize - FStream.Size) = 0 then
          Exit(True);

        while (totalSize - FStream.Size) > 0 do
        begin
          size := totalSize - FStream.Size;
          if size = 0 then
            Exit(True)
          else if size < 0 then
          begin
            FError := 'file size error';
            Exit;
          end
          else if size > FDownloadChunkSize then
            size := FDownloadChunkSize;

          if not DownloadChunk(FStream, url, FStream.Size, FStream.Size + size - 1) then
            Exit;
        end;

        if (totalSize - FStream.Size) = 0 then
          Exit(True);
      except
        on E: Exception do
          FError := E.Message;
      end;
    finally
      if Assigned(FStream) then
        FreeAndNil(FStream);
    end;
  end;
end;

function TFileService.GetStatusCode: Integer;
begin
  if Assigned(FResponse) then
    Result := FResponse.StatusCode;
end;

function TFileService.GetUrlDownload: string;
begin
  Result := Format(URL_DOWNLOAD, [FHost, FFileName])
end;

function TFileService.GetUrlUpload: string;
begin
  Result := Format(URL_UPLOAD, [FHost, FFileName])
end;

function TFileService.Info(const fileName: string; fileInfoStream: TMemoryStream; md5: Boolean = False): Boolean;
var
  url: string;
begin
  Result := False;

  if md5 then
    url := URL_INFO_MD5
  else
    url := URL_INFO;

  url := Format(url, [FHost, fileName]);
  try
    FResponse := Self.Get(url, fileInfoStream, nil);
    if FResponse.StatusCode = 200 then
      Exit(True);
  except
    on E: Exception do
      FError := E.Message;
  end;
end;

function TFileService.InfoHead(const fileName: string; var Size: Int64): Boolean;
var
  url: string;
begin
  Result := False;

  url := Format(URL_INFO_HEAD, [FHost, fileName]);
  try
    FResponse := Self.Head(url, nil);
    if FResponse.StatusCode = 200 then
    begin
      Size := FResponse.ContentLength;
      Exit(True);
    end
    else
      FError := 'server file not found';
  except
    on E: Exception do
      FError := E.Message;
  end;
end;

procedure TFileService.SetWorkDir(const dir: string);
begin
  if not DirectoryExists(dir) then
    MkDir(dir);
  FWorkDir := dir;
end;

function TFileService.UploadChunk(fileStream: TFileStream; const url: string; const offset, size: Int64): Boolean;
var
  AHeaders: TNetHeaders;
  SStream: TStringStream;
  MStream, RespStream: TMemoryStream;
begin
  Result := False;

  SetLength(AHeaders, Length(AHeaders) + 1);
  AHeaders[High(AHeaders)] := TNameValuePair.Create('Range', Format(RANGE_BYTES, [offset, size]));

  MStream := TMemoryStream.Create;
  RespStream := TMemoryStream.Create;
  try
    try
      fileStream.Position := offset;

      MStream.CopyFrom(fileStream, size - offset + 1);
      MStream.Position := 0;

      FResponse := Self.Post(url, MStream, RespStream, AHeaders);
      if (FResponse.StatusCode = 200) or (FResponse.StatusCode = 206) then
      begin
        Exit(True);
      end
      else
      begin
        SStream := TStringStream.Create;
        try
          SStream.LoadFromStream(RespStream);
          FError := SStream.DataString;
        finally
          if Assigned(SStream) then
            FreeAndNil(SStream);
        end;
      end;
    except
      on E: Exception do
        FError := E.Message
    end;
  finally
    if Assigned(MStream) then
      FreeAndNil(MStream);
    if Assigned(RespStream) then
      FreeAndNil(RespStream);
  end;
end;

function TFileService.UploadFile(const fileName: string): Boolean;
var
  fStream: TFileStream;
  sSize, size: Int64;
  url: string;
begin
  Result := False;
  if not FileExists(fileName) then
  begin
    FError := 'file not found';
    Exit;
  end;

  sSize := 0;
  size := 0;

  if InfoHead(fileName, sSize) or (StatusCode = 404) then
  begin
    url := Format(URL_UPLOAD, [FHost, fileName]);
    try
      try
        fStream := TFileStream.Create(fileName, fmOpenRead or fmShareExclusive);

        if (sSize - fStream.Size) = 0 then
          Exit(True);

        while (fStream.Size - sSize) > 0 do
        begin
          size := fStream.Size - sSize;
          if size = 0 then
            Exit(True)
          else if size < 0 then
          begin
            FError := 'server file size error';
            Exit;
          end
          else if size > FUploadChunkSize then
            size := FUploadChunkSize;

          if not UploadChunk(fStream, url, sSize, sSize + size - 1) then
            Exit;

          sSize := sSize + size;
        end;

        if (sSize - fStream.Size) = 0 then
          Exit(True);
      except
        on E: Exception do
          FError := E.Message;
      end;
    finally
      if Assigned(fStream) then
        FreeAndNil(fStream);
    end;
  end;
end;

end.

