import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { taskApi, type TaskStatus } from '../services/taskApi';
import { toDataUri } from '../utils/imageUtils';
import {
  MdArrowBack,
  MdCheckCircle,
  MdError,
  MdSchedule,
  MdCancel,
  MdAccessTime,
  MdPerson,
  MdTheaters,
  MdImage,
  MdClose,
} from 'react-icons/md';

const TaskDetailPage: React.FC = () => {
  const { taskId } = useParams<{ taskId: string }>();
  const navigate = useNavigate();
  const [taskData, setTaskData] = useState<TaskStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedImage, setSelectedImage] = useState<string | null>(null);

  // 轮询任务状态
  useEffect(() => {
    if (!taskId) return;

    const fetchTaskStatus = async () => {
      try {
        const response = await taskApi.getTaskStatus(taskId);
        setTaskData(response.data);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : '获取任务状态失败');
      } finally {
        setLoading(false);
      }
    };

    // 立即获取一次
    fetchTaskStatus();

    // 设置轮询定时器 (每2秒)
    const pollInterval = setInterval(() => {
      if (taskData?.status === 'processing' || taskData?.status === 'pending') {
        fetchTaskStatus();
      } else {
        clearInterval(pollInterval);
      }
    }, 2000);

    // 清理定时器
    return () => clearInterval(pollInterval);
  }, [taskId, taskData?.status]);

  const statusConfig = {
    pending: {
      label: '等待中',
      bgColor: 'bg-gray-100',
      textColor: 'text-gray-700',
      icon: <MdSchedule className="w-5 h-5" />,
      ringColor: 'ring-gray-300',
    },
    processing: {
      label: '处理中',
      bgColor: 'bg-primary-50',
      textColor: 'text-primary-700',
      icon: <MdAccessTime className="w-5 h-5 animate-spin" />,
      ringColor: 'ring-primary-300',
    },
    completed: {
      label: '已完成',
      bgColor: 'bg-success-50',
      textColor: 'text-success-700',
      icon: <MdCheckCircle className="w-5 h-5" />,
      ringColor: 'ring-success-300',
    },
    failed: {
      label: '失败',
      bgColor: 'bg-danger-50',
      textColor: 'text-danger-700',
      icon: <MdError className="w-5 h-5" />,
      ringColor: 'ring-danger-300',
    },
    cancelled: {
      label: '已取消',
      bgColor: 'bg-warning-50',
      textColor: 'text-warning-700',
      icon: <MdCancel className="w-5 h-5" />,
      ringColor: 'ring-warning-300',
    },
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="text-center animate-fade-in">
          <div className="relative w-16 h-16 mx-auto mb-4">
            <div className="absolute inset-0 border-4 border-primary-200 rounded-full"></div>
            <div className="absolute inset-0 border-4 border-primary-500 rounded-full border-t-transparent animate-spin"></div>
          </div>
          <p className="text-gray-600 font-medium">加载任务详情...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center p-4">
        <div className="bg-white rounded-2xl shadow-lg border border-danger-200 p-8 max-w-md w-full animate-scale-in">
          <div className="flex items-center justify-center w-16 h-16 bg-danger-100 rounded-full mx-auto mb-4">
            <MdError className="w-8 h-8 text-danger-600" />
          </div>
          <h3 className="text-xl font-semibold text-gray-900 text-center mb-2">加载失败</h3>
          <p className="text-gray-600 text-center mb-6">{error}</p>
          <button
            onClick={() => navigate('/tasks')}
            className="w-full px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 shadow-md hover:shadow-lg transition-all duration-200 font-medium"
          >
            返回任务列表
          </button>
        </div>
      </div>
    );
  }

  if (!taskData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600">任务不存在</p>
        </div>
      </div>
    );
  }

  const config = statusConfig[taskData.status];

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
        {/* Back Button */}
        <button
          onClick={() => navigate('/tasks')}
          className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-6 transition-colors group"
        >
          <MdArrowBack className="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
          <span className="font-medium">返回任务列表</span>
        </button>

        {/* Header Card */}
        <div className="bg-white rounded-2xl shadow-md border border-gray-200 p-6 sm:p-8 mb-6 animate-slide-down">
          <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-4">
            <div className="flex-1">
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
                {taskData.task_id ? `任务详情` : '任务'}
              </h1>
              <p className="text-sm text-gray-500 font-mono">ID: {taskData.task_id}</p>
            </div>
            <span className={`inline-flex items-center gap-2 px-4 py-2.5 rounded-xl text-sm font-semibold ${config.bgColor} ${config.textColor} ring-2 ${config.ringColor} whitespace-nowrap`}>
              {config.icon}
              {config.label}
            </span>
          </div>
        </div>

        {/* Progress Card */}
        <div className="bg-white rounded-2xl shadow-md border border-gray-200 p-6 sm:p-8 mb-6 animate-slide-up" style={{ animationDelay: '100ms' }}>
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center">
              <MdAccessTime className="w-5 h-5 text-primary-600" />
            </div>
            <h2 className="text-xl font-semibold text-gray-900">执行进度</h2>
          </div>

          <div className="space-y-4">
            <div>
              <div className="flex justify-between items-baseline mb-3">
                <span className="text-gray-700 font-medium">{taskData.progress.current_step}</span>
                <span className="text-2xl font-bold text-primary-600">{taskData.progress.percentage}%</span>
              </div>
              <div className="relative w-full bg-gray-200 rounded-full h-3 overflow-hidden">
                <div
                  className="absolute inset-y-0 left-0 bg-gradient-to-r from-primary-500 to-primary-600 rounded-full transition-all duration-500 ease-out"
                  style={{ width: `${taskData.progress.percentage}%` }}
                />
                {taskData.status === 'processing' && (
                  <div
                    className="absolute inset-y-0 left-0 bg-primary-400 rounded-full animate-pulse"
                    style={{ width: `${taskData.progress.percentage}%` }}
                  />
                )}
              </div>
            </div>

            <div className="flex items-center justify-between text-sm text-gray-600 pt-2">
              <span>步骤进度</span>
              <span className="font-semibold">
                {taskData.progress.current_step_index} / {taskData.progress.total_steps}
              </span>
            </div>

            {taskData.progress.details && (
              <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-6 pt-6 border-t border-gray-100">
                <div className="bg-gradient-to-br from-primary-50 to-primary-100 p-4 rounded-xl border border-primary-200">
                  <div className="flex items-center gap-2 mb-2">
                    <MdPerson className="w-4 h-4 text-primary-600" />
                    <p className="text-xs text-primary-600 font-medium">角色提取</p>
                  </div>
                  <p className="text-2xl font-bold text-primary-700">{taskData.progress.details.characters_extracted}</p>
                </div>
                <div className="bg-gradient-to-br from-success-50 to-success-100 p-4 rounded-xl border border-success-200">
                  <div className="flex items-center gap-2 mb-2">
                    <MdImage className="w-4 h-4 text-success-600" />
                    <p className="text-xs text-success-600 font-medium">参考图生成</p>
                  </div>
                  <p className="text-2xl font-bold text-success-700">{taskData.progress.details.characters_generated}</p>
                </div>
                <div className="bg-gradient-to-br from-warning-50 to-warning-100 p-4 rounded-xl border border-warning-200">
                  <div className="flex items-center gap-2 mb-2">
                    <MdTheaters className="w-4 h-4 text-warning-600" />
                    <p className="text-xs text-warning-600 font-medium">场景划分</p>
                  </div>
                  <p className="text-2xl font-bold text-warning-700">{taskData.progress.details.scenes_divided}</p>
                </div>
                <div className="bg-gradient-to-br from-purple-50 to-purple-100 p-4 rounded-xl border border-purple-200">
                  <div className="flex items-center gap-2 mb-2">
                    <MdImage className="w-4 h-4 text-purple-600" />
                    <p className="text-xs text-purple-600 font-medium">场景生成</p>
                  </div>
                  <p className="text-2xl font-bold text-purple-700">{taskData.progress.details.scenes_generated}</p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Result Card */}
        {taskData.result && (
          <div className="bg-white rounded-2xl shadow-md border border-gray-200 p-6 sm:p-8 mb-6 animate-slide-up" style={{ animationDelay: '200ms' }}>
            <div className="flex items-center gap-3 mb-6">
              <div className="w-10 h-10 bg-success-100 rounded-lg flex items-center justify-center">
                <MdCheckCircle className="w-5 h-5 text-success-600" />
              </div>
              <h2 className="text-xl font-semibold text-gray-900">生成结果</h2>
            </div>

            {/* Summary Stats */}
            {taskData.result && (
              <div className="grid grid-cols-2 gap-6 mb-8">
                <div className="bg-gradient-to-br from-primary-50 to-primary-100 p-6 rounded-xl border border-primary-200">
                  <div className="flex items-center gap-3 mb-2">
                    <MdPerson className="w-6 h-6 text-primary-600" />
                    <p className="text-gray-700 font-medium">角色数量</p>
                  </div>
                  <p className="text-4xl font-bold text-primary-700">{taskData.result.character_count}</p>
                </div>
                <div className="bg-gradient-to-br from-success-50 to-success-100 p-6 rounded-xl border border-success-200">
                  <div className="flex items-center gap-3 mb-2">
                    <MdTheaters className="w-6 h-6 text-success-600" />
                    <p className="text-gray-700 font-medium">场景数量</p>
                  </div>
                  <p className="text-4xl font-bold text-success-700">{taskData.result.scene_count}</p>
                </div>
              </div>
            )}

            {/* Characters */}
            {taskData.result && taskData.result.characters.length > 0 && (
              <div className="mb-8">
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
                  <MdPerson className="w-5 h-5 text-primary-600" />
                  角色列表
                </h3>
                <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
                  {taskData.result.characters.map((char) => (
                    <div
                      key={char.id}
                      className="group bg-gradient-to-br from-gray-50 to-white border border-gray-200 rounded-xl overflow-hidden hover:shadow-lg transition-all duration-300 cursor-pointer"
                      onClick={() => char.reference_image_url && setSelectedImage(toDataUri(char.reference_image_url))}
                    >
                      {char.reference_image_url && (
                        <div className="relative aspect-[3/4] overflow-hidden bg-gray-100">
                          <img
                            src={toDataUri(char.reference_image_url)}
                            alt={char.name}
                            className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                          />
                          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                        </div>
                      )}
                      <div className="p-3">
                        <p className="font-semibold text-center text-gray-900 truncate">{char.name}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Scenes */}
            {taskData.result && taskData.result.scenes.length > 0 && (
              <div>
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
                  <MdTheaters className="w-5 h-5 text-success-600" />
                  场景列表
                </h3>
                <div className="space-y-4">
                  {taskData.result.scenes.map((scene) => (
                    <div
                      key={scene.id}
                      className="group bg-gradient-to-br from-gray-50 to-white border border-gray-200 rounded-xl p-5 hover:shadow-lg transition-all duration-300"
                    >
                      <div className="flex items-start gap-4">
                        <div className="flex-shrink-0 w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center">
                          <span className="text-primary-700 font-bold text-lg">{scene.sequence_num}</span>
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-gray-700 mb-3 leading-relaxed">{scene.description}</p>
                          {scene.image_url && (
                            <div
                              className="relative rounded-lg overflow-hidden bg-gray-100 cursor-pointer"
                              onClick={() => setSelectedImage(toDataUri(scene.image_url!))}
                            >
                              <img
                                src={toDataUri(scene.image_url)}
                                alt={`Scene ${scene.sequence_num}`}
                                className="w-full h-auto object-contain group-hover:scale-105 transition-transform duration-300"
                              />
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        {/* Error Card */}
        {taskData.error && (
          <div className="bg-white rounded-2xl shadow-md border-2 border-danger-200 p-6 sm:p-8 mb-6 animate-slide-up">
            <div className="flex items-start gap-4">
              <div className="flex-shrink-0 w-12 h-12 bg-danger-100 rounded-lg flex items-center justify-center">
                <MdError className="w-6 h-6 text-danger-600" />
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-danger-800 mb-2">错误信息</h3>
                <p className="text-danger-700 mb-2">{taskData.error.message}</p>
                <p className="text-sm text-danger-600">错误代码: {taskData.error.code}</p>
              </div>
            </div>
          </div>
        )}

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-4 animate-fade-in">
          <button
            onClick={() => navigate('/tasks')}
            className="flex-1 sm:flex-initial px-6 py-3 bg-gray-600 text-white rounded-xl hover:bg-gray-700 shadow-md hover:shadow-lg transition-all duration-200 font-medium"
          >
            返回任务列表
          </button>
          {(taskData.status === 'pending' || taskData.status === 'processing') && (
            <button
              onClick={async () => {
                if (confirm('确定要取消任务吗？')) {
                  try {
                    await taskApi.cancelTask(taskId!);
                    setTaskData({ ...taskData, status: 'cancelled' });
                  } catch (err) {
                    alert('取消任务失败');
                  }
                }
              }}
              className="flex-1 sm:flex-initial px-6 py-3 bg-danger-500 text-white rounded-xl hover:bg-danger-600 shadow-md hover:shadow-lg transition-all duration-200 font-medium"
            >
              取消任务
            </button>
          )}
        </div>
      </div>

      {/* Image Modal */}
      {selectedImage && (
        <div
          className="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-fade-in"
          onClick={() => setSelectedImage(null)}
        >
          <div className="relative max-w-5xl w-full animate-scale-in">
            <button
              onClick={() => setSelectedImage(null)}
              className="absolute -top-12 right-0 p-2 bg-white/10 hover:bg-white/20 rounded-lg text-white transition-colors"
            >
              <MdClose className="w-6 h-6" />
            </button>
            <img
              src={selectedImage}
              alt="预览"
              className="w-full h-auto rounded-xl shadow-2xl"
              onClick={(e) => e.stopPropagation()}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default TaskDetailPage;
