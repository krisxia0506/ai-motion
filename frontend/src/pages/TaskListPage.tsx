import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { taskApi, type TaskListItem, type TaskStatusType } from '../services/taskApi';
import { MdAdd, MdRefresh, MdCheckCircle, MdError, MdSchedule, MdCancel, MdAccessTime } from 'react-icons/md';

const TaskListPage: React.FC = () => {
  const navigate = useNavigate();
  const [tasks, setTasks] = useState<TaskListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [statusFilter, setStatusFilter] = useState<string>('');
  const pageSize = 12;

  // 获取任务列表
  const fetchTasks = async () => {
    try {
      setLoading(true);
      const response = await taskApi.getTaskList({
        page: currentPage,
        page_size: pageSize,
        status: statusFilter || undefined
      });
      setTasks(response.data.items);
      setTotalPages(response.data.pagination.total_pages);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取任务列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, [currentPage, statusFilter]);

  // 状态配置 - 使用更现代的配置
  const statusConfig: Record<TaskStatusType, {
    label: string;
    bgColor: string;
    textColor: string;
    icon: React.ReactNode;
    ringColor: string;
  }> = {
    pending: {
      label: '等待中',
      bgColor: 'bg-gray-100',
      textColor: 'text-gray-700',
      icon: <MdSchedule className="w-4 h-4" />,
      ringColor: 'ring-gray-200'
    },
    processing: {
      label: '处理中',
      bgColor: 'bg-primary-50',
      textColor: 'text-primary-700',
      icon: <MdAccessTime className="w-4 h-4 animate-spin" />,
      ringColor: 'ring-primary-200'
    },
    completed: {
      label: '已完成',
      bgColor: 'bg-success-50',
      textColor: 'text-success-700',
      icon: <MdCheckCircle className="w-4 h-4" />,
      ringColor: 'ring-success-200'
    },
    failed: {
      label: '失败',
      bgColor: 'bg-danger-50',
      textColor: 'text-danger-700',
      icon: <MdError className="w-4 h-4" />,
      ringColor: 'ring-danger-200'
    },
    cancelled: {
      label: '已取消',
      bgColor: 'bg-warning-50',
      textColor: 'text-warning-700',
      icon: <MdCancel className="w-4 h-4" />,
      ringColor: 'ring-warning-200'
    },
  };

  const formatDate = (date: string) => {
    const d = new Date(date);
    const now = new Date();
    const diffMs = now.getTime() - d.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return '刚刚';
    if (diffMins < 60) return `${diffMins}分钟前`;
    if (diffHours < 24) return `${diffHours}小时前`;
    if (diffDays < 7) return `${diffDays}天前`;

    return d.toLocaleString('zh-CN', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading && tasks.length === 0) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="text-center animate-fade-in">
          <div className="relative w-16 h-16 mx-auto mb-4">
            <div className="absolute inset-0 border-4 border-primary-200 rounded-full"></div>
            <div className="absolute inset-0 border-4 border-primary-500 rounded-full border-t-transparent animate-spin"></div>
          </div>
          <p className="text-gray-600 font-medium">加载任务列表...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
        {/* Header with improved spacing and hierarchy */}
        <div className="mb-8 sm:mb-10 animate-slide-down">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-4">
            <div>
              <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-2">我的任务</h1>
              <p className="text-gray-600">管理和追踪您的漫画生成任务</p>
            </div>
            <button
              onClick={() => navigate('/create')}
              className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-gradient-to-r from-primary-500 to-primary-600 text-white rounded-xl hover:from-primary-600 hover:to-primary-700 shadow-lg hover:shadow-xl transition-all duration-200 transform hover:scale-105 font-medium"
            >
              <MdAdd className="w-5 h-5" />
              <span>创建新任务</span>
            </button>
          </div>

          {/* Filter Bar with improved design */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-4 sm:p-5">
            <div className="flex flex-col sm:flex-row sm:items-center gap-4">
              <label className="text-sm font-medium text-gray-700 whitespace-nowrap">筛选条件</label>
              <div className="flex flex-wrap items-center gap-3 flex-1">
                <select
                  value={statusFilter}
                  onChange={(e) => {
                    setStatusFilter(e.target.value);
                    setCurrentPage(1);
                  }}
                  className="flex-1 sm:flex-initial min-w-[200px] px-4 py-2.5 bg-gray-50 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 text-sm font-medium"
                >
                  <option value="">全部状态</option>
                  <option value="completed">已完成</option>
                  <option value="processing">处理中</option>
                  <option value="failed">失败</option>
                  <option value="pending">等待中</option>
                  <option value="cancelled">已取消</option>
                </select>
                <button
                  onClick={fetchTasks}
                  disabled={loading}
                  className="inline-flex items-center gap-2 px-4 py-2.5 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed font-medium text-sm"
                >
                  <MdRefresh className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                  <span className="hidden sm:inline">刷新</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Error Alert with better design */}
        {error && (
          <div className="mb-6 bg-danger-50 border-l-4 border-danger-500 rounded-lg p-4 animate-slide-up">
            <div className="flex items-start gap-3">
              <MdError className="w-5 h-5 text-danger-600 flex-shrink-0 mt-0.5" />
              <div>
                <h3 className="text-sm font-semibold text-danger-800">加载失败</h3>
                <p className="text-sm text-danger-700 mt-1">{error}</p>
              </div>
            </div>
          </div>
        )}

        {/* Task Grid or Empty State */}
        {tasks.length === 0 ? (
          <div className="bg-white rounded-2xl shadow-sm border border-gray-200 p-12 sm:p-16 text-center animate-fade-in">
            <div className="max-w-md mx-auto">
              <div className="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <MdSchedule className="w-10 h-10 text-gray-400" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">暂无任务</h3>
              <p className="text-gray-600 mb-6">开始创建您的第一个漫画生成任务吧</p>
              <button
                onClick={() => navigate('/create')}
                className="inline-flex items-center gap-2 px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 shadow-md hover:shadow-lg transition-all duration-200 font-medium"
              >
                <MdAdd className="w-5 h-5" />
                创建第一个任务
              </button>
            </div>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
              {tasks.map((task, index) => {
                const config = statusConfig[task.status];
                return (
                  <div
                    key={task.task_id}
                    onClick={() => navigate(`/task/${task.task_id}`)}
                    className="group bg-white rounded-xl shadow-card hover:shadow-card-hover border border-gray-200 hover:border-primary-300 p-6 cursor-pointer transition-all duration-300 transform hover:-translate-y-1 animate-slide-up"
                    style={{ animationDelay: `${index * 50}ms` }}
                  >
                    {/* Header: Title and Status Badge */}
                    <div className="flex items-start justify-between gap-3 mb-4">
                      <h3 className="text-lg font-semibold text-gray-900 flex-1 line-clamp-2 group-hover:text-primary-600 transition-colors">
                        {task.title || '未命名任务'}
                      </h3>
                      <span className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold ${config.bgColor} ${config.textColor} ring-1 ${config.ringColor} whitespace-nowrap`}>
                        {config.icon}
                        {config.label}
                      </span>
                    </div>

                    {/* Progress Bar */}
                    <div className="mb-4">
                      <div className="flex justify-between items-baseline text-sm mb-2">
                        <span className="text-gray-600 font-medium">{task.progress.current_step}</span>
                        <span className="text-primary-600 font-bold">{task.progress.percentage}%</span>
                      </div>
                      <div className="relative w-full bg-gray-200 rounded-full h-2.5 overflow-hidden">
                        <div
                          className="absolute inset-y-0 left-0 bg-gradient-to-r from-primary-500 to-primary-600 rounded-full transition-all duration-500 ease-out"
                          style={{ width: `${task.progress.percentage}%` }}
                        />
                        {task.status === 'processing' && (
                          <div
                            className="absolute inset-y-0 left-0 bg-primary-400 rounded-full animate-pulse"
                            style={{ width: `${task.progress.percentage}%` }}
                          />
                        )}
                      </div>
                    </div>

                    {/* Stats for completed tasks */}
                    {task.status === 'completed' && ((task.character_count ?? 0) > 0 || (task.scene_count ?? 0) > 0) && (
                      <div className="grid grid-cols-2 gap-3 mb-4">
                        <div className="bg-gradient-to-br from-primary-50 to-primary-100 p-3 rounded-lg border border-primary-200">
                          <p className="text-xs text-primary-600 font-medium mb-1">角色</p>
                          <p className="text-2xl font-bold text-primary-700">{task.character_count ?? 0}</p>
                        </div>
                        <div className="bg-gradient-to-br from-success-50 to-success-100 p-3 rounded-lg border border-success-200">
                          <p className="text-xs text-success-600 font-medium mb-1">场景</p>
                          <p className="text-2xl font-bold text-success-700">{task.scene_count ?? 0}</p>
                        </div>
                      </div>
                    )}

                    {/* Error for failed tasks */}
                    {task.status === 'failed' && task.error && (
                      <div className="bg-danger-50 border border-danger-200 rounded-lg p-3 mb-4">
                        <div className="flex items-start gap-2">
                          <MdError className="w-4 h-4 text-danger-600 flex-shrink-0 mt-0.5" />
                          <p className="text-sm text-danger-700 line-clamp-2 font-medium">{task.error.message}</p>
                        </div>
                      </div>
                    )}

                    {/* Timestamp with icon */}
                    <div className="flex items-center gap-2 text-sm text-gray-500 pt-4 border-t border-gray-100">
                      <MdAccessTime className="w-4 h-4" />
                      {task.status === 'completed' && task.completed_at ? (
                        <span>完成于 {formatDate(task.completed_at)}</span>
                      ) : task.status === 'failed' && task.failed_at ? (
                        <span>失败于 {formatDate(task.failed_at)}</span>
                      ) : (
                        <span>创建于 {formatDate(task.created_at)}</span>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Improved Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-2 animate-fade-in">
                <button
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                  className="px-4 py-2.5 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 hover:border-gray-400 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white transition-all duration-200"
                >
                  上一页
                </button>

                <div className="flex items-center gap-1">
                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    let pageNum;
                    if (totalPages <= 5) {
                      pageNum = i + 1;
                    } else if (currentPage <= 3) {
                      pageNum = i + 1;
                    } else if (currentPage >= totalPages - 2) {
                      pageNum = totalPages - 4 + i;
                    } else {
                      pageNum = currentPage - 2 + i;
                    }
                    return (
                      <button
                        key={pageNum}
                        onClick={() => setCurrentPage(pageNum)}
                        className={`min-w-[40px] px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200 ${
                          currentPage === pageNum
                            ? 'bg-primary-500 text-white shadow-md scale-105'
                            : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 hover:border-gray-400'
                        }`}
                      >
                        {pageNum}
                      </button>
                    );
                  })}
                </div>

                <button
                  onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                  disabled={currentPage === totalPages}
                  className="px-4 py-2.5 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 hover:border-gray-400 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-white transition-all duration-200"
                >
                  下一页
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default TaskListPage;
