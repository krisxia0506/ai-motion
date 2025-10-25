import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { MdUploadFile, MdAutoAwesome, MdPerson, MdMovie, MdVolumeUp } from 'react-icons/md';
import { Button, Card, CardBody } from '../components/common';
import { novelApi } from '../services/novelApi';
import './HomePage.css';

function HomePage() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      setError(null);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    try {
      setUploading(true);
      setError(null);

      await novelApi.uploadNovelFile(selectedFile);
      navigate('/novels');
    } catch {
      setError('上传失败，请重试');
    } finally {
      setUploading(false);
    }
  };

  const features = [
    {
      icon: <MdAutoAwesome size={40} />,
      title: '智能解析',
      description: '自动解析小说内容，提取章节、场景和对话',
    },
    {
      icon: <MdPerson size={40} />,
      title: '角色识别',
      description: '智能识别角色，生成一致性角色形象',
    },
    {
      icon: <MdMovie size={40} />,
      title: '场景生成',
      description: '基于AI生成高质量动漫场景图片',
    },
    {
      icon: <MdVolumeUp size={40} />,
      title: 'AI配音',
      description: '自动生成角色配音，支持多种音色',
    },
  ];

  return (
    <div className="home-page">
      <div className="container">
        <section className="hero-section">
          <div className="hero-content">
            <h1 className="hero-title">
              <span className="hero-title-gradient">AI-Motion</span>
              <br />
              智能动漫生成系统
            </h1>
            <p className="hero-description">
              将文字小说转化为精美动漫，AI驱动的创作新体验
            </p>
          </div>

          <Card className="upload-card">
            <CardBody>
              <h2 className="upload-title">
                <MdUploadFile size={28} />
                上传小说开始创作
              </h2>

              <div className="upload-area">
                <input
                  type="file"
                  accept=".txt,.epub,.pdf"
                  onChange={handleFileChange}
                  id="file-upload"
                  className="file-input"
                />
                <label htmlFor="file-upload" className="file-label">
                  {selectedFile ? (
                    <div className="file-selected">
                      <MdUploadFile size={48} />
                      <span className="file-name">{selectedFile.name}</span>
                      <span className="file-size">
                        {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                      </span>
                    </div>
                  ) : (
                    <div className="file-placeholder">
                      <MdUploadFile size={48} />
                      <span>点击选择文件或拖拽文件到此处</span>
                      <span className="file-hint">支持 .txt, .epub, .pdf 格式</span>
                    </div>
                  )}
                </label>
              </div>

              {error && <div className="upload-error">{error}</div>}

              <div className="upload-actions">
                {selectedFile && (
                  <Button
                    variant="outline"
                    onClick={() => {
                      setSelectedFile(null);
                      setError(null);
                    }}
                  >
                    重新选择
                  </Button>
                )}
                <Button
                  variant="primary"
                  size="large"
                  fullWidth={!selectedFile}
                  onClick={handleUpload}
                  disabled={!selectedFile || uploading}
                  loading={uploading}
                >
                  {uploading ? '上传中...' : '开始上传'}
                </Button>
              </div>
            </CardBody>
          </Card>
        </section>

        <section className="features-section">
          <h2 className="section-title">核心功能</h2>
          <div className="features-grid">
            {features.map((feature, index) => (
              <Card key={index} hoverable className="feature-card">
                <CardBody>
                  <div className="feature-icon">{feature.icon}</div>
                  <h3 className="feature-title">{feature.title}</h3>
                  <p className="feature-description">{feature.description}</p>
                </CardBody>
              </Card>
            ))}
          </div>
        </section>

        <section className="workflow-section">
          <h2 className="section-title">工作流程</h2>
          <div className="workflow-steps">
            {[
              '上传小说文件',
              '自动解析内容',
              '识别角色与场景',
              '生成动漫画面',
              '添加AI配音',
              '导出完整视频',
            ].map((step, index) => (
              <div key={index} className="workflow-step">
                <div className="step-number">{index + 1}</div>
                <div className="step-text">{step}</div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}

export default HomePage;
