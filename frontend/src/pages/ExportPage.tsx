import { MdFileDownload } from 'react-icons/md';

function ExportPage() {
  return (
    <div className="container" style={{ padding: '48px 0', textAlign: 'center' }}>
      <MdFileDownload size={64} style={{ color: 'var(--color-primary)', margin: '0 auto 24px' }} />
      <h1>导出视频</h1>
      <p style={{ color: 'var(--color-text-secondary)' }}>此功能正在开发中...</p>
    </div>
  );
}

export default ExportPage;
