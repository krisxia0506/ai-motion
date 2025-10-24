import { MdMovie } from 'react-icons/md';

function GenerationPage() {
  return (
    <div className="container" style={{ padding: '48px 0', textAlign: 'center' }}>
      <MdMovie size={64} style={{ color: 'var(--color-primary)', margin: '0 auto 24px' }} />
      <h1>生成动漫</h1>
      <p style={{ color: 'var(--color-text-secondary)' }}>此功能正在开发中...</p>
    </div>
  );
}

export default GenerationPage;
