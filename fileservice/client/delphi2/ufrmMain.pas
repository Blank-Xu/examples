unit ufrmMain;

interface

uses
  System.SysUtils,
  System.Types,
  System.UITypes,
  System.Classes,
  System.Variants,
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
  ufmeUploadFile,
  System.Net.URLClient,
  System.Net.HttpClient,
  System.Net.HttpClientComponent;

type
  TfrmMain = class(TForm)
    pnl1: TPanel;
    edtHost: TEdit;
    edtFileName: TEdit;
    btnInfoHead: TButton;
    btnInfo: TButton;
    btnDownload: TButton;
    btnUpload: TButton;
    btnDelete: TButton;
    mmo: TMemo;
    Http1: TNetHTTPClient;
    pnlPB: TPanel;
    procedure btnDownloadClick(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure btnUploadClick(Sender: TObject);
    procedure btnDeleteClick(Sender: TObject);
  private
    const
      FWorkDir = 'files';
  public
		{ Public declarations }
  end;

var
  frmMain: TfrmMain;

implementation

{$R *.fmx}

procedure TfrmMain.btnDeleteClick(Sender: TObject);
var
  fme: TFileService;
  msg: string;
begin
  fme := TFileService.Create(nil, edtHost.Text, edtFileName.Text);
  try
    if fme.DeleteFile then
      msg := 'delete success'
    else
      msg := 'delete failed, err: ' + fme.Error;

    mmo.Lines.Add(msg);
  finally
    fme.Free;
  end;
end;

procedure TfrmMain.btnDownloadClick(Sender: TObject);
begin
  TThread.CreateAnonymousThread(
    procedure
    var
      fmeDownload: TfmeDownloadFile;
      Msg: string;
    begin
      fmeDownload := TfmeDownloadFile.Create(pnlPB, edtHost.Text, edtFileName.Text, FWorkDir);
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

procedure TfrmMain.btnUploadClick(Sender: TObject);
begin
  TThread.CreateAnonymousThread(
    procedure
    var
      fmeUpload: TfmeUploadFile;
      Msg: string;
    begin
      fmeUpload := TfmeUploadFile.Create(pnlPB, edtHost.Text, edtFileName.Text, FWorkDir);
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

procedure TfrmMain.FormCreate(Sender: TObject);
begin
  if not DirectoryExists(FWorkDir) then
    MkDir(FWorkDir);
end;

end.

