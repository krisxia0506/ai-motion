import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { taskApi } from '../services/taskApi';
import {
  MdArrowBack,
  MdUploadFile,
  MdTitle,
  MdPerson,
  MdDescription,
  MdCheckCircle,
  MdError,
  MdLightbulb,
} from 'react-icons/md';

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

  const contentLength = content.length;
  const titleLength = title.length;
  const isContentValid = contentLength >= 100 && contentLength <= 5000;
  const isTitleValid = titleLength > 0 && titleLength <= 200;
  const canSubmit = isTitleValid && isContentValid && !loading;

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
        {/* Back Button */}
        <button
          onClick={() => navigate('/tasks')}
          className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-6 transition-colors group"
        >
          <MdArrowBack className="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
          <span className="font-medium">返回任务列表</span>
        </button>

        {/* Header */}
        <div className="mb-8 animate-slide-down">
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-2">创建漫画生成任务</h1>
          <p className="text-gray-600">输入小说内容，AI 将自动提取角色、划分场景并生成精美漫画</p>
        </div>

        <div className="grid lg:grid-cols-3 gap-6">
          {/* Main Form */}
          <div className="lg:col-span-2">
            <form onSubmit={handleSubmit} className="bg-white rounded-2xl shadow-md border border-gray-200 p-6 sm:p-8 animate-slide-up">
              {/* Error Alert */}
              {error && (
                <div className="mb-6 bg-danger-50 border-l-4 border-danger-500 rounded-lg p-4 animate-slide-down">
                  <div className="flex items-start gap-3">
                    <MdError className="w-5 h-5 text-danger-600 flex-shrink-0 mt-0.5" />
                    <div>
                      <h3 className="text-sm font-semibold text-danger-800">输入错误</h3>
                      <p className="text-sm text-danger-700 mt-1">{error}</p>
                    </div>
                  </div>
                </div>
              )}

              {/* Title Field */}
              <div className="mb-6">
                <label className="flex items-center gap-2 text-gray-900 font-semibold mb-3">
                  <MdTitle className="w-5 h-5 text-primary-600" />
                  <span>标题</span>
                  <span className="text-danger-500">*</span>
                </label>
                <input
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="请输入小说标题，例如：三体、哈利波特..."
                  maxLength={200}
                  className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 ${
                    titleLength > 0 && !isTitleValid ? 'border-danger-300 bg-danger-50' : 'border-gray-200'
                  }`}
                  disabled={loading}
                />
                <div className="flex items-center justify-between mt-2">
                  <p className={`text-sm font-medium ${
                    titleLength === 0 ? 'text-gray-500' :
                    isTitleValid ? 'text-success-600' : 'text-danger-600'
                  }`}>
                    {titleLength > 0 && (
                      isTitleValid ? (
                        <span className="inline-flex items-center gap-1">
                          <MdCheckCircle className="w-4 h-4" />
                          标题格式正确
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1">
                          <MdError className="w-4 h-4" />
                          标题不能为空
                        </span>
                      )
                    )}
                  </p>
                  <p className="text-sm text-gray-500">{titleLength}/200</p>
                </div>
              </div>

              {/* Author Field */}
              <div className="mb-6">
                <label className="flex items-center gap-2 text-gray-900 font-semibold mb-3">
                  <MdPerson className="w-5 h-5 text-primary-600" />
                  <span>作者</span>
                  <span className="text-gray-400 text-sm font-normal">(可选)</span>
                </label>
                <input
                  type="text"
                  value={author}
                  onChange={(e) => setAuthor(e.target.value)}
                  placeholder="请输入作者名，例如：刘慈欣、J.K.罗琳..."
                  className="w-full px-4 py-3 bg-gray-50 border-2 border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200"
                  disabled={loading}
                />
              </div>

              {/* Content Field */}
              <div className="mb-6">
                <div className="flex items-center justify-between mb-3">
                  <label className="flex items-center gap-2 text-gray-900 font-semibold">
                    <MdDescription className="w-5 h-5 text-primary-600" />
                    <span>内容</span>
                    <span className="text-danger-500">*</span>
                  </label>
                  <label className="inline-flex items-center gap-2 px-4 py-2 bg-primary-50 text-primary-700 rounded-lg hover:bg-primary-100 cursor-pointer transition-colors duration-200 font-medium text-sm border border-primary-200">
                    <MdUploadFile className="w-4 h-4" />
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
                  placeholder="请输入或粘贴小说内容...&#10;&#10;支持完整章节或片段，内容应包含清晰的角色描述和故事情节。&#10;&#10;示例：&#10;在一个阳光明媚的午后，年轻的魔法师艾莉丝站在古老的魔法学院门前。她有着一头乌黑的长发和明亮的蓝色眼睛..."
                  rows={12}
                  maxLength={5000}
                  className={`w-full px-4 py-3 bg-gray-50 border-2 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 font-mono text-sm resize-y ${
                    contentLength > 0 && !isContentValid ? 'border-danger-300 bg-danger-50' : 'border-gray-200'
                  }`}
                  disabled={loading}
                />
                <div className="flex items-center justify-between mt-2">
                  <p className={`text-sm font-medium ${
                    contentLength === 0 ? 'text-gray-500' :
                    contentLength < 100 ? 'text-warning-600' :
                    isContentValid ? 'text-success-600' : 'text-danger-600'
                  }`}>
                    {contentLength === 0 ? (
                      '至少需要 100 个字符'
                    ) : contentLength < 100 ? (
                      <span className="inline-flex items-center gap-1">
                        <MdError className="w-4 h-4" />
                        还需要 {100 - contentLength} 个字符
                      </span>
                    ) : isContentValid ? (
                      <span className="inline-flex items-center gap-1">
                        <MdCheckCircle className="w-4 h-4" />
                        内容格式正确
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1">
                        <MdError className="w-4 h-4" />
                        超出字符限制
                      </span>
                    )}
                  </p>
                  <p className={`text-sm font-medium ${
                    contentLength > 5000 ? 'text-danger-600' : 'text-gray-500'
                  }`}>
                    {contentLength}/5000
                  </p>
                </div>
              </div>

              {/* Actions */}
              <div className="flex flex-col sm:flex-row gap-4 pt-4 border-t border-gray-200">
                <button
                  type="button"
                  onClick={() => navigate('/tasks')}
                  className="flex-1 sm:flex-initial px-6 py-3 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition-colors duration-200 font-medium"
                  disabled={loading}
                >
                  取消
                </button>
                <button
                  type="submit"
                  className="flex-1 px-8 py-3 bg-gradient-to-r from-primary-500 to-primary-600 text-white rounded-xl hover:from-primary-600 hover:to-primary-700 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:shadow-lg transition-all duration-200 font-semibold"
                  disabled={!canSubmit}
                >
                  {loading ? (
                    <span className="inline-flex items-center gap-2">
                      <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                      </svg>
                      创建中...
                    </span>
                  ) : (
                    '创建任务'
                  )}
                </button>
              </div>
            </form>
          </div>

          {/* Tips Sidebar */}
          <div className="lg:col-span-1 space-y-6 animate-slide-up" style={{ animationDelay: '100ms' }}>
            {/* Quick Tips */}
            <div className="bg-gradient-to-br from-primary-50 to-primary-100 border border-primary-200 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center gap-3 mb-4">
                <div className="w-10 h-10 bg-primary-500 rounded-lg flex items-center justify-center">
                  <MdLightbulb className="w-5 h-5 text-white" />
                </div>
                <h3 className="font-semibold text-primary-900">创作提示</h3>
              </div>
              <ul className="space-y-3 text-sm text-primary-800">
                <li className="flex items-start gap-2">
                  <span className="text-primary-600 font-bold mt-0.5">•</span>
                  <span>内容应包含<strong>完整的故事情节</strong>和<strong>清晰的角色描述</strong></span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-primary-600 font-bold mt-0.5">•</span>
                  <span>描述角色的<strong>外貌特征</strong>有助于生成更准确的形象</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-primary-600 font-bold mt-0.5">•</span>
                  <span>包含<strong>对话和场景描写</strong>可以生成更丰富的漫画内容</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-primary-600 font-bold mt-0.5">•</span>
                  <span>支持导入 <code className="px-1.5 py-0.5 bg-primary-200 rounded text-xs">.txt</code> 文件</span>
                </li>
              </ul>
            </div>

            {/* Process Info */}
            <div className="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm">
              <h3 className="font-semibold text-gray-900 mb-4">生成流程</h3>
              <ol className="space-y-4">
                <li className="flex gap-3">
                  <div className="flex-shrink-0 w-8 h-8 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center font-bold text-sm">
                    1
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">角色提取</p>
                    <p className="text-sm text-gray-600 mt-1">AI 自动识别和提取故事中的角色</p>
                  </div>
                </li>
                <li className="flex gap-3">
                  <div className="flex-shrink-0 w-8 h-8 bg-success-100 text-success-700 rounded-full flex items-center justify-center font-bold text-sm">
                    2
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">生成参考图</p>
                    <p className="text-sm text-gray-600 mt-1">为每个角色生成一致的参考形象</p>
                  </div>
                </li>
                <li className="flex gap-3">
                  <div className="flex-shrink-0 w-8 h-8 bg-warning-100 text-warning-700 rounded-full flex items-center justify-center font-bold text-sm">
                    3
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">场景划分</p>
                    <p className="text-sm text-gray-600 mt-1">智能划分故事场景和对话</p>
                  </div>
                </li>
                <li className="flex gap-3">
                  <div className="flex-shrink-0 w-8 h-8 bg-purple-100 text-purple-700 rounded-full flex items-center justify-center font-bold text-sm">
                    4
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">生成漫画</p>
                    <p className="text-sm text-gray-600 mt-1">为每个场景生成精美的漫画图片</p>
                  </div>
                </li>
              </ol>
            </div>

            {/* Estimated Time */}
            <div className="bg-gradient-to-br from-gray-50 to-white border border-gray-200 rounded-2xl p-6 shadow-sm">
              <h3 className="font-semibold text-gray-900 mb-3">预计时间</h3>
              <div className="flex items-center gap-3">
                <div className="flex-shrink-0 w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center">
                  <svg className="w-6 h-6 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div>
                  <p className="text-2xl font-bold text-gray-900">3-5 分钟</p>
                  <p className="text-sm text-gray-600">根据内容长度而定</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateMangaPage;
