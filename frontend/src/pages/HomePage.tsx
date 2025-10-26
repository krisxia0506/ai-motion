import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { MdUploadFile, MdAutoAwesome, MdPerson, MdMovie, MdVolumeUp, MdEdit } from 'react-icons/md';
import { Button, Card, CardBody } from '../components/common';
import { apiClient } from '../services/api';
import './HomePage.css';

type InputMode = 'file' | 'text';

function HomePage() {
  const [inputMode, setInputMode] = useState<InputMode>('file');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [textContent, setTextContent] = useState('');
  const [title, setTitle] = useState('');
  const [author, setAuthor] = useState('');
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      setError(null);
      if (!title) {
        setTitle(file.name.replace(/\.[^/.]+$/, ''));
      }
    }
  };

  const readFileContent = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target?.result as string;
        resolve(content);
      };
      reader.onerror = () => reject(new Error('Failed to read file'));
      reader.readAsText(file);
    });
  };

  const handleGenerate = async () => {
    if (!title) {
      setError('请输入标题');
      return;
    }

    if (inputMode === 'file' && !selectedFile) {
      setError('请选择文件');
      return;
    }

    if (inputMode === 'text' && !textContent.trim()) {
      setError('请输入小说内容');
      return;
    }

    try {
      setUploading(true);
      setError(null);

      let content = '';
      if (inputMode === 'file' && selectedFile) {
        content = await readFileContent(selectedFile);
      } else {
        content = textContent;
      }

      const response = await apiClient.post<{ novel_id: string }>('/manga/generate', {
        title,
        author: author || 'Unknown',
        content,
      });

      navigate(`/novels/${response.data.novel_id}`);
    } catch (err) {
      const error = err instanceof Error ? err : new Error('生成失败');
      setError(error.message);
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
                <MdAutoAwesome size={28} />
                生成漫画开始创作
              </h2>

              {/* Input Mode Switcher */}
              <div className="mode-switcher">
                <button
                  type="button"
                  className={`mode-btn ${inputMode === 'file' ? 'active' : ''}`}
                  onClick={() => setInputMode('file')}
                  disabled={uploading}
                >
                  <MdUploadFile size={20} />
                  <span>上传文件</span>
                </button>
                <button
                  type="button"
                  className={`mode-btn ${inputMode === 'text' ? 'active' : ''}`}
                  onClick={() => setInputMode('text')}
                  disabled={uploading}
                >
                  <MdEdit size={20} />
                  <span>输入文本</span>
                </button>
              </div>

              {/* File Upload Mode */}
              {inputMode === 'file' && (
                <div className="upload-area">
                  <input
                    type="file"
                    accept=".txt,.doc,.docx"
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
                        <span className="file-hint">支持 .txt, .doc, .docx 格式</span>
                      </div>
                    )}
                  </label>
                </div>
              )}

              {/* Text Input Mode */}
              {inputMode === 'text' && (
                <div className="text-input-area">
                  <textarea
                    value={textContent}
                    onChange={(e) => setTextContent(e.target.value)}
                    placeholder="在此粘贴或输入小说内容..."
                    disabled={uploading}
                    className="text-input"
                    rows={6}
                  />
                </div>
              )}

              {/* Title and Author Fields */}
              <div className="metadata-fields">
                <div className="form-field">
                  <label htmlFor="title" className="field-label">
                    标题 *
                  </label>
                  <input
                    id="title"
                    type="text"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="请输入漫画标题"
                    disabled={uploading}
                    className="field-input"
                  />
                </div>
                <div className="form-field">
                  <label htmlFor="author" className="field-label">
                    作者
                  </label>
                  <input
                    id="author"
                    type="text"
                    value={author}
                    onChange={(e) => setAuthor(e.target.value)}
                    placeholder="请输入作者名称（可选）"
                    disabled={uploading}
                    className="field-input"
                  />
                </div>
              </div>

              {error && <div className="upload-error">{error}</div>}

              <div className="upload-actions">
                <Button
                  variant="primary"
                  size="large"
                  fullWidth
                  onClick={handleGenerate}
                  disabled={
                    !title ||
                    uploading ||
                    (inputMode === 'file' && !selectedFile) ||
                    (inputMode === 'text' && !textContent.trim())
                  }
                  loading={uploading}
                >
                  {uploading ? '生成中...' : '开始生成'}
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
              '上传文件或输入文本',
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
