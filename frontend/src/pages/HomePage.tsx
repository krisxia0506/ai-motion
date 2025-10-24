import { useState } from 'react';

function HomePage() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    // TODO: 实现文件上传逻辑
    console.log('Uploading file:', selectedFile.name);
  };

  return (
    <div className="container">
      <h1>AI-Motion - 智能动漫生成系统</h1>

      <div className="upload-section">
        <h2>上传小说</h2>
        <input
          type="file"
          accept=".txt,.epub,.pdf"
          onChange={handleFileChange}
        />
        {selectedFile && (
          <div>
            <p>已选择文件: {selectedFile.name}</p>
            <button onClick={handleUpload}>开始上传</button>
          </div>
        )}
      </div>

      <div className="features">
        <h2>主要功能</h2>
        <ul>
          <li>自动解析小说内容</li>
          <li>智能识别角色</li>
          <li>生成角色一致性图片</li>
          <li>自动配音</li>
          <li>导出动漫视频</li>
        </ul>
      </div>
    </div>
  );
}

export default HomePage;
