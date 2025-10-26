import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { taskApi, TaskListItem, TaskStatus } from '../services/taskApi';

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
      const response = await taskApi.getTaskList(currentPage, pageSize, statusFilter);
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

  const statusColors: Record<TaskStatus, string> = {
    pending: 'bg-gray-200 text-gray-800',
    processing: 'bg-blue-500 text-white',
    completed: 'bg-green-500 text-white',
    failed: 'bg-red-500 text-white',
    cancelled: 'bg-orange-500 text-white',
  };

  const statusLabels: Record<TaskStatus, string> = {
    pending: '等待中',
    processing: '处理中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消',
  };

  const formatDate = (date: string) => {
    return new Date(date).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading && tasks.length === 0) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-3xl font-bold">我的任务</h1>
          <button
            onClick={() => navigate('/create')}
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
          >
            创建新任务
          </button>
        </div>

        {/* Filter */}
        <div className="bg-white rounded-lg shadow-md p-4 mb-6">
          <div className="flex items-center gap-4">
            <label className="text-gray-700 font-medium">状态筛选:</label>
            <select
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value);
                setCurrentPage(1);
              }}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">全部</option>
              <option value="completed">已完成</option>
              <option value="processing">处理中</option>
              <option value="failed">失败</option>
              <option value="pending">等待中</option>
              <option value="cancelled">已取消</option>
            </select>
            <button
              onClick={fetchTasks}
              className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
            >
              刷新
            </button>
          </div>
        </div>

        {/* Error */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <p className="text-red-600">{error}</p>
          </div>
        )}

        {/* Task Grid */}
        {tasks.length === 0 ? (
          <div className="bg-white rounded-lg shadow-md p-12 text-center">
            <p className="text-gray-500 text-lg">暂无任务</p>
            <button
              onClick={() => navigate('/create')}
              className="mt-4 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              创建第一个任务
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-6">
            {tasks.map((task) => (
              <div
                key={task.task_id}
                onClick={() => navigate(`/task/${task.task_id}`)}
                className="bg-white rounded-lg shadow-md p-6 cursor-pointer hover:shadow-lg transition-shadow"
              >
                {/* Title and Status */}
                <div className="flex items-start justify-between mb-4">
                  <h3 className="text-lg font-semibold flex-1 mr-2 line-clamp-2">
                    {task.title || '未命名任务'}
                  </h3>
                  <span className={`px-3 py-1 rounded-full text-sm whitespace-nowrap ${statusColors[task.status]}`}>
                    {statusLabels[task.status]}
                  </span>
                </div>

                {/* Progress */}
                <div className="mb-4">
                  <div className="flex justify-between text-sm text-gray-600 mb-1">
                    <span>{task.progress.current_step}</span>
                    <span>{task.progress.percentage}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                      style={{ width: `${task.progress.percentage}%` }}
                    />
                  </div>
                </div>

                {/* Stats for completed tasks */}
                {task.status === 'completed' && (
                  <div className="grid grid-cols-2 gap-3 mb-4">
                    <div className="bg-blue-50 p-3 rounded">
                      <p className="text-xs text-gray-600">角色</p>
                      <p className="text-xl font-bold text-blue-600">{task.character_count || 0}</p>
                    </div>
                    <div className="bg-green-50 p-3 rounded">
                      <p className="text-xs text-gray-600">场景</p>
                      <p className="text-xl font-bold text-green-600">{task.scene_count || 0}</p>
                    </div>
                  </div>
                )}

                {/* Error for failed tasks */}
                {task.status === 'failed' && task.error && (
                  <div className="bg-red-50 border border-red-200 rounded p-3 mb-4">
                    <p className="text-sm text-red-600 line-clamp-2">{task.error.message}</p>
                  </div>
                )}

                {/* Timestamp */}
                <div className="text-sm text-gray-500">
                  {task.status === 'completed' && task.completed_at ? (
                    <span>完成于 {formatDate(task.completed_at)}</span>
                  ) : task.status === 'failed' && task.failed_at ? (
                    <span>失败于 {formatDate(task.failed_at)}</span>
                  ) : (
                    <span>创建于 {formatDate(task.created_at)}</span>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-center gap-2">
            <button
              onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
              disabled={currentPage === 1}
              className="px-4 py-2 bg-white border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
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
                    className={`px-4 py-2 rounded-lg ${
                      currentPage === pageNum
                        ? 'bg-blue-600 text-white'
                        : 'bg-white border border-gray-300 hover:bg-gray-50'
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
              className="px-4 py-2 bg-white border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              下一页
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default TaskListPage;
