import { MdPerson } from 'react-icons/md';

function CharacterPage() {
  return (
    <div className="container" style={{ padding: '48px 0', textAlign: 'center' }}>
      <MdPerson size={64} style={{ color: 'var(--color-primary)', margin: '0 auto 24px' }} />
      <h1>角色管理</h1>
      <p style={{ color: 'var(--color-text-secondary)' }}>此功能正在开发中...</p>
    </div>
  );
}

export default CharacterPage;
