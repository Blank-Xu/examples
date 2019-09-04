unit uProgressFileStream;

// https://codeday.me/bug/20190519/1133458.html

interface

uses
  System.Classes,
  System.SysUtils;

type
  TProgressFileStreamOnProgress = procedure(Sender: TObject; Processed: Int64; Size: Int64; ContentLength: Int64; TimeStart: cardinal) of object;

  TProgressFileStream = class(TFileStream)
  private
    FProcessed: Int64;
    FContentLength: Int64;
    FTimeStart: cardinal;
    FBytesDiff: cardinal;
    FSize: Int64;
    FOnProgress: TProgressFileStreamOnProgress;
    procedure Init;
    procedure DoProgress(const AProcessed: Longint);
  protected
    procedure SetSize(NewSize: Longint); overload; override;
  public
    constructor Create(const AFileName: string; Mode: Word); overload;
    constructor Create(const AFileName: string; Mode: Word; Rights: Cardinal); overload;
    function Read(var Buffer; Count: Longint): Longint; overload; override;
    function Write(const Buffer; Count: Longint): Longint; overload; override;
    function Read(Buffer: TBytes; Offset, Count: Longint): Longint; overload; override;
    function Write(const Buffer: TBytes; Offset, Count: Longint): Longint; overload; override;
    function Seek(const Offset: Int64; Origin: TSeekOrigin): Int64; overload; override;
    property OnProgress: TProgressFileStreamOnProgress read FOnProgress write FOnProgress;
    property ContentLength: Int64 read FContentLength write FContentLength;
    property TimeStart: cardinal read FTimeStart write FTimeStart;
    property BytesDiff: cardinal read FBytesDiff write FBytesDiff;
  end;

implementation

{ TProgressFileStream }

constructor TProgressFileStream.Create(const AFileName: string; Mode: Word);
begin
  inherited Create(AFileName, Mode);

  Init;
end;

constructor TProgressFileStream.Create(const AFileName: string; Mode: Word; Rights: Cardinal);
begin
  inherited Create(AFileName, Mode, Rights);

  Init;
end;

function TProgressFileStream.Read(var Buffer; Count: Longint): Longint;
begin
  Result := inherited Read(Buffer, Count);

  DoProgress(Result);
end;

function TProgressFileStream.Write(const Buffer; Count: Longint): Longint;
begin
  Result := inherited Write(Buffer, Count);

  DoProgress(Result);
end;

function TProgressFileStream.Read(Buffer: TBytes; Offset, Count: Longint): Longint;
begin
  Result := inherited Read(Buffer, Offset, Count);

  DoProgress(Result);
end;

function TProgressFileStream.Write(const Buffer: TBytes; Offset, Count: Longint): Longint;
begin
  Result := inherited Write(Buffer, Offset, Count);

  DoProgress(Result);
end;

function TProgressFileStream.Seek(const Offset: Int64; Origin: TSeekOrigin): Int64;
begin
  Result := inherited Seek(Offset, Origin);

  if Origin <> soCurrent then
    FProcessed := Result;
end;

procedure TProgressFileStream.SetSize(NewSize: Longint);
begin
  inherited SetSize(NewSize);

  FSize := NewSize;
end;

procedure TProgressFileStream.Init;
const
  BYTES_DIFF = 1024 * 100;
begin
  FOnProgress := nil;
  FProcessed := 0;
  FContentLength := 0;
  FTimeStart := TThread.GetTickCount;
  FBytesDiff := BYTES_DIFF;
  FSize := Size;
end;

procedure TProgressFileStream.DoProgress(const AProcessed: Longint);
var
  aCurrentProcessed: Longint;
begin
  if not (Assigned(FOnProgress)) then
    Exit;

  aCurrentProcessed := FProcessed;

  Inc(FProcessed, AProcessed);

  if FContentLength = 0 then
    FContentLength := FSize;

  if (FProcessed = FSize) or (FBytesDiff = 0) or (aCurrentProcessed - FBytesDiff < FProcessed) then
    FOnProgress(Self, FProcessed, FSize, FContentLength, FTimeStart);
end;

end.

