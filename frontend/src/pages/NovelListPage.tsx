import React from 'react';
import { useNavigate } from 'react-router-dom';
import { MdAdd, MdBook, MdCheckCircle, MdError } from 'react-icons/md';
import { Button, Card, CardBody, LoadingSpinner } from '../components/common';
import './NovelListPage.css';

function NovelListPage() {
  const navigate = useNavigate();
  const [loading] = React.useState(false);

  const mockNovels = [
    { id: '1', title: '示例小说 1', author: '作者A', status: 'completed' },
    { id: '2', title: '示例小说 2', author: '作者B', status: 'parsing' },
  ];

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <MdCheckCircle className="status-icon status-completed" />;
      case 'parsing':
        return <LoadingSpinner size="small" />;
      default:
        return <MdError className="status-icon status-error" />;
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'completed':
        return '已完成';
      case 'parsing':
        return '解析中';
      default:
        return '失败';
    }
  };

  if (loading) {
    return <LoadingSpinner fullScreen text="加载中..." />;
  }

  return (
    <div className="novel-list-page">
      <div className="container">
        <div className="page-header">
          <h1>我的小说</h1>
          <Button
            variant="primary"
            onClick={() => navigate('/')}
          >
            <MdAdd size={20} />
            上传新小说
          </Button>
        </div>

        <div className="novels-grid">
          {mockNovels.map((novel) => (
            <Card
              key={novel.id}
              hoverable
              onClick={() => navigate(`/novels/${novel.id}`)}
              className="novel-card"
            >
              <CardBody>
                <div className="novel-icon">
                  <MdBook size={48} />
                </div>
                <h3 className="novel-title">{novel.title}</h3>
                <p className="novel-author">作者: {novel.author}</p>
                <div className="novel-status">
                  {getStatusIcon(novel.status)}
                  <span>{getStatusText(novel.status)}</span>
                </div>
              </CardBody>
            </Card>
          ))}
        </div>

        {mockNovels.length === 0 && (
          <div className="empty-state">
            <MdBook size={64} />
            <h2>还没有小说</h2>
            <p>开始上传你的第一本小说吧！</p>
            <Button variant="primary" onClick={() => navigate('/')}>
              上传小说
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

export default NovelListPage;
