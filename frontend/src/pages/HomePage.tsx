import { useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { taskApi } from '../services/taskApi';
import {
  MdArrowForward,
  MdAutoAwesome,
  MdPerson,
  MdMovie,
  MdImage,
  MdLightbulb,
  MdCheckCircle,
  MdSend,
  MdError
} from 'react-icons/md';

function HomePage() {
  const navigate = useNavigate();
  const [content, setContent] = useState('在一个阳光明媚的午后，年轻的魔法师艾莉丝站在古老的魔法学院门前。她有着一头乌黑的长发和明亮的蓝色眼睛，穿着一身深蓝色的魔法师长袍。今天是她成为正式魔法师的第一天，心中既兴奋又紧张。学院的大门缓缓打开，里面传来悠扬的魔法钟声。');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const features = [
    {
      icon: <MdAutoAwesome className="w-10 h-10" />,
      title: '智能解析',
      description: '自动解析小说内容，提取章节、场景和对话',
      color: 'from-purple-500 to-purple-600',
      bgColor: 'bg-purple-50',
      iconColor: 'text-purple-600'
    },
    {
      icon: <MdPerson className="w-10 h-10" />,
      title: '角色识别',
      description: '智能识别角色，生成一致性角色形象',
      color: 'from-primary-500 to-primary-600',
      bgColor: 'bg-primary-50',
      iconColor: 'text-primary-600'
    },
    {
      icon: <MdMovie className="w-10 h-10" />,
      title: '场景生成',
      description: '基于AI生成高质量动漫场景图片',
      color: 'from-success-500 to-success-600',
      bgColor: 'bg-success-50',
      iconColor: 'text-success-600'
    },
    {
      icon: <MdImage className="w-10 h-10" />,
      title: '图像渲染',
      description: '精美的图像渲染和视觉效果处理',
      color: 'from-warning-500 to-warning-600',
      bgColor: 'bg-warning-50',
      iconColor: 'text-warning-600'
    },
  ];

  const workflowSteps = [
    { title: '上传内容', desc: '输入小说文本或上传文件' },
    { title: '智能解析', desc: 'AI 自动分析内容结构' },
    { title: '角色提取', desc: '识别并生成角色形象' },
    { title: '场景划分', desc: '智能划分故事场景' },
    { title: '图像生成', desc: '生成精美漫画画面' },
    { title: '完成导出', desc: '获取完整的漫画作品' },
  ];

  const highlights = [
    '支持多种文件格式',
    'AI 驱动的智能处理',
    '角色形象一致性保证',
    '高质量图像生成',
    '快速处理速度',
    '简单易用的界面'
  ];

  const handleQuickGenerate = async (e: React.FormEvent) => {
    e.preventDefault();

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
        title: content.substring(0, 50) + '...',
        author: '匿名作者',
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

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-white to-gray-100">
      {/* Hero Section */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-primary-500/10 via-purple-500/10 to-success-500/10 pointer-events-none" />

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 sm:py-28">
          <div className="text-center mb-12 animate-slide-down">
            <div className="inline-flex items-center gap-2 px-4 py-2 bg-primary-100 text-primary-700 rounded-full text-sm font-semibold mb-6">
              <MdAutoAwesome className="w-4 h-4" />
              <span>AI 驱动的智能创作平台</span>
            </div>

            <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold text-gray-900 mb-6">
              <span className="bg-gradient-to-r from-primary-600 via-purple-600 to-success-600 bg-clip-text text-transparent">
                AI-Motion
              </span>
              <br />
              <span className="text-gray-800">智能动漫生成系统</span>
            </h1>

            <p className="text-xl sm:text-2xl text-gray-600 mb-8 max-w-3xl mx-auto">
              将文字小说转化为精美动漫，体验 AI 驱动的创作新时代
            </p>

            {/* Quick Generate Form */}
            <div className="max-w-4xl mx-auto mb-8">
              <form onSubmit={handleQuickGenerate} className="bg-white rounded-2xl shadow-xl border-2 border-gray-200 p-6">
                {error && (
                  <div className="mb-4 bg-danger-50 border-l-4 border-danger-500 rounded-lg p-4 animate-slide-down">
                    <div className="flex items-start gap-3">
                      <MdError className="w-5 h-5 text-danger-600 flex-shrink-0 mt-0.5" />
                      <p className="text-sm text-danger-700">{error}</p>
                    </div>
                  </div>
                )}

                <div className="space-y-4">
                  <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="输入小说内容，至少100个字符...&#10;&#10;示例：在一个阳光明媚的午后，年轻的魔法师艾莉丝站在古老的魔法学院门前。她有着一头乌黑的长发和明亮的蓝色眼睛，穿着一身深蓝色的魔法师长袍..."
                    rows={6}
                    maxLength={5000}
                    className="w-full px-4 py-3 bg-gray-50 border-2 border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 resize-none"
                    disabled={loading}
                  />

                  <div className="flex items-center justify-between">
                    <span className={`text-sm font-medium ${
                      content.length === 0 ? 'text-gray-500' :
                      content.length < 100 ? 'text-warning-600' :
                      content.length <= 5000 ? 'text-success-600' : 'text-danger-600'
                    }`}>
                      {content.length === 0 ? '至少需要 100 个字符' :
                       content.length < 100 ? `还需要 ${100 - content.length} 个字符` :
                       content.length <= 5000 ? '内容格式正确' : '超出字符限制'}
                    </span>
                    <span className="text-sm text-gray-500">{content.length}/5000</span>
                  </div>

                  <button
                    type="submit"
                    disabled={loading || content.length < 100 || content.length > 5000}
                    className="w-full group inline-flex items-center justify-center gap-3 px-8 py-4 bg-gradient-to-r from-primary-500 to-primary-600 text-white rounded-xl hover:from-primary-600 hover:to-primary-700 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:shadow-lg transition-all duration-300 font-semibold text-lg"
                  >
                    {loading ? (
                      <>
                        <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <span>生成中...</span>
                      </>
                    ) : (
                      <>
                        <MdSend className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                        <span>立即生成动漫</span>
                      </>
                    )}
                  </button>
                </div>
              </form>
            </div>

            <div className="flex items-center justify-center">
              <button
                onClick={() => navigate('/tasks')}
                className="inline-flex items-center gap-2 px-8 py-4 bg-white text-gray-700 rounded-2xl hover:bg-gray-50 border-2 border-gray-200 hover:border-gray-300 shadow-lg hover:shadow-xl transition-all duration-300 font-semibold text-lg"
              >
                查看任务
              </button>
            </div>
          </div>

          {/* Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6 mt-16 animate-slide-up">
            {[
              { number: '5分钟', label: '平均生成时间' },
              { number: '高质量', label: 'AI 图像生成' },
              { number: '自动化', label: '智能处理流程' },
              { number: '易使用', label: '简单友好界面' },
            ].map((stat, index) => (
              <div
                key={index}
                className="bg-white rounded-2xl p-6 text-center shadow-card hover:shadow-card-hover transition-all duration-300 border border-gray-200"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <p className="text-3xl font-bold text-gray-900 mb-2">{stat.number}</p>
                <p className="text-sm text-gray-600">{stat.label}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 sm:py-28 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16 animate-slide-down">
            <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4">核心功能</h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              强大的 AI 技术，为您的创作赋能
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature, index) => (
              <div
                key={index}
                className="group bg-gradient-to-br from-gray-50 to-white rounded-2xl p-8 border border-gray-200 hover:border-gray-300 shadow-card hover:shadow-card-hover transition-all duration-300 transform hover:-translate-y-2 animate-slide-up"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <div className={`inline-flex items-center justify-center w-16 h-16 ${feature.bgColor} rounded-2xl mb-6 group-hover:scale-110 transition-transform duration-300`}>
                  <div className={feature.iconColor}>
                    {feature.icon}
                  </div>
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-3">{feature.title}</h3>
                <p className="text-gray-600 leading-relaxed">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Workflow Section */}
      <section className="py-20 sm:py-28 bg-gradient-to-br from-gray-50 to-gray-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16 animate-slide-down">
            <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4">工作流程</h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              简单六步，完成您的动漫创作
            </p>
          </div>

          <div className="relative">
            {/* Connection Line */}
            <div className="hidden lg:block absolute top-1/2 left-0 right-0 h-1 bg-gradient-to-r from-primary-200 via-purple-200 to-success-200 -translate-y-1/2 -z-10" />

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-6 gap-8">
              {workflowSteps.map((step, index) => (
                <div
                  key={index}
                  className="relative animate-slide-up"
                  style={{ animationDelay: `${index * 100}ms` }}
                >
                  <div className="bg-white rounded-2xl p-6 shadow-card hover:shadow-card-hover transition-all duration-300 transform hover:-translate-y-2 border border-gray-200">
                    <div className="flex flex-col items-center text-center">
                      <div className="w-12 h-12 bg-gradient-to-br from-primary-500 to-primary-600 text-white rounded-full flex items-center justify-center font-bold text-lg mb-4 shadow-lg">
                        {index + 1}
                      </div>
                      <h3 className="font-bold text-gray-900 mb-2">{step.title}</h3>
                      <p className="text-sm text-gray-600">{step.desc}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Highlights Section */}
      <section className="py-20 sm:py-28 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="bg-gradient-to-br from-primary-50 to-purple-50 rounded-3xl p-12 sm:p-16 border border-primary-200 animate-slide-up">
            <div className="grid lg:grid-cols-2 gap-12 items-center">
              <div>
                <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-6">
                  为什么选择<br />
                  <span className="bg-gradient-to-r from-primary-600 to-purple-600 bg-clip-text text-transparent">
                    AI-Motion？
                  </span>
                </h2>
                <p className="text-xl text-gray-600 mb-8">
                  我们提供业界领先的 AI 动漫生成技术，让创作变得简单而高效
                </p>
                <button
                  onClick={() => navigate('/create')}
                  className="group inline-flex items-center gap-3 px-8 py-4 bg-gradient-to-r from-primary-500 to-primary-600 text-white rounded-2xl hover:from-primary-600 hover:to-primary-700 shadow-lg hover:shadow-xl transition-all duration-300 font-semibold"
                >
                  <span>立即开始</span>
                  <MdArrowForward className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                </button>
              </div>

              <div className="space-y-4">
                {highlights.map((highlight, index) => (
                  <div
                    key={index}
                    className="flex items-center gap-4 bg-white rounded-xl p-4 shadow-sm hover:shadow-md transition-all duration-300 animate-slide-up"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <div className="flex-shrink-0 w-10 h-10 bg-success-100 rounded-lg flex items-center justify-center">
                      <MdCheckCircle className="w-6 h-6 text-success-600" />
                    </div>
                    <span className="text-gray-900 font-medium">{highlight}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 sm:py-28 bg-gradient-to-br from-primary-600 to-purple-600 text-white">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center animate-slide-up">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-white/20 rounded-2xl mb-6 backdrop-blur-sm">
            <MdLightbulb className="w-8 h-8" />
          </div>
          <h2 className="text-4xl sm:text-5xl font-bold mb-6">
            准备好开始创作了吗？
          </h2>
          <p className="text-xl text-white/90 mb-10 max-w-2xl mx-auto">
            只需几分钟，就能将您的故事转化为精美的动漫作品
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <button
              onClick={() => navigate('/create')}
              className="group inline-flex items-center gap-3 px-8 py-4 bg-white text-primary-600 rounded-2xl hover:bg-gray-50 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:scale-105 font-semibold text-lg"
            >
              <span>立即创建任务</span>
              <MdArrowForward className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </button>
            <button
              onClick={() => navigate('/tasks')}
              className="inline-flex items-center gap-2 px-8 py-4 bg-white/10 text-white rounded-2xl hover:bg-white/20 backdrop-blur-sm border-2 border-white/30 shadow-lg hover:shadow-xl transition-all duration-300 font-semibold text-lg"
            >
              查看示例
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}

export default HomePage;
