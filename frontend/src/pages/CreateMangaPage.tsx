import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { taskApi } from '../services/taskApi';

const CreateMangaPage: React.FC = () => {
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [author, setAuthor] = useState('');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!title.trim()) {
      setError('请输入标题');
      return;
    }

    if (content.length < 100) {
      setError('内容至少需要100个字符');
      return;
    }

    if (content.length > 5000) {
      setError('内容不能超过5000个字符');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await taskApi.createTask({
        title,
        author: author || undefined,
        content,
      });

      // 导航到任务详情页
      navigate(`/task/${response.data.task_id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : '创建任务失败');
    } finally {
      setLoading(false);
    }
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith('.txt')) {
      setError('只支持 .txt 文件');
      return;
    }

    try {
      const text = await file.text();
      setContent(text);
      setError(null);
    } catch (err) {
      setError('读取文件失败');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold mb-2">创建漫画生成任务</h1>
          <p className="text-gray-600">输入小说内容，AI将自动生成角色和场景</p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow-md p-6">
          {/* Error */}
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
              <p className="text-red-600">{error}</p>
            </div>
          )}

          {/* Title */}
          <div className="mb-6">
            <label className="block text-gray-700 font-medium mb-2">
              标题 <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="请输入小说标题"
              maxLength={200}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
            <p className="text-sm text-gray-500 mt-1">
              {title.length}/200 字符
            </p>
          </div>

          {/* Author */}
          <div className="mb-6">
            <label className="block text-gray-700 font-medium mb-2">作者</label>
            <input
              type="text"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              placeholder="请输入作者名（可选）"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>

          {/* Content */}
          <div className="mb-6">
            <div className="flex items-center justify-between mb-2">
              <label className="block text-gray-700 font-medium">
                内容 <span className="text-red-500">*</span>
              </label>
              <label className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 cursor-pointer">
                <input
                  type="file"
                  accept=".txt"
                  onChange={handleFileUpload}
                  className="hidden"
                  disabled={loading}
                />
                从文件导入
              </label>
            </div>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="请输入或粘贴小说内容（至少100字符，最多5000字符）"
              rows={15}
              maxLength={5000}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm resize-y"
              disabled={loading}
            />
            <p className="text-sm text-gray-500 mt-1">
              {content.length}/5000 字符
              {content.length < 100 && content.length > 0 && (
                <span className="text-red-500 ml-2">还需要 {100 - content.length} 个字符</span>
              )}
            </p>
          </div>

          {/* Actions */}
          <div className="flex items-center justify-between">
            <button
              type="button"
              onClick={() => navigate('/tasks')}
              className="px-6 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
              disabled={loading}
            >
              取消
            </button>
            <button
              type="submit"
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed font-medium"
              disabled={loading || !title.trim() || content.length < 100}
            >
              {loading ? '创建中...' : '创建任务'}
            </button>
          </div>
        </form>

        {/* Tips */}
        <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-6">
          <h3 className="font-semibold text-blue-800 mb-2">提示</h3>
          <ul className="text-sm text-blue-700 space-y-1">
            <li>• 内容应包含完整的故事情节和角色描述</li>
            <li>• AI将自动提取角色、划分场景并生成图片</li>
            <li>• 生成过程需要几分钟，您可以在任务详情页查看进度</li>
            <li>• 支持导入 .txt 文件快速填充内容</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default CreateMangaPage;
