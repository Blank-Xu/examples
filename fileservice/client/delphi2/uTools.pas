unit uTools;

interface

uses
  System.Classes,
  System.SysUtils;

const
	KB = 1024;
  MB = 1024 * KB;

function BytesToStr(ABytes: Integer): string;

implementation

function BytesToStr(ABytes: Integer): string;
var
	iKb: Integer;
begin
	iKb := Round(ABytes / KB);
	if iKb > 1000 then
		Result := Format('%.2f MB/s', [iKb / KB])
	else
		Result := Format('%d KB/s', [iKb]);
end;

end.

