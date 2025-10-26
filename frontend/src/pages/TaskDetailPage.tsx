import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { taskApi, type TaskStatus } from '../services/taskApi';

const TaskDetailPage: React.FC = () => {
  const { taskId } = useParams<{ taskId: string }>();
  const navigate = useNavigate();
  const [taskData, setTaskData] = useState<TaskStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">加载中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 max-w-md">
          <h3 className="text-red-800 font-semibold">错误</h3>
          <p className="text-red-600 mt-2">{error}</p>
          <button
            onClick={() => navigate('/tasks')}
            className="mt-4 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
          >
            返回任务列表
          </button>
        </div>
      </div>
    );
  }

  if (!taskData) {
    return <div>任务不存在</div>;
  }

  const statusColors = {
    pending: 'bg-gray-200 text-gray-800',
    processing: 'bg-blue-500 text-white',
    completed: 'bg-green-500 text-white',
    failed: 'bg-red-500 text-white',
    cancelled: 'bg-orange-500 text-white',
  };

  const statusLabels = {
    pending: '等待中',
    processing: '处理中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消',
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">任务详情</h1>
            <span className={`px-4 py-2 rounded-full ${statusColors[taskData.status]}`}>
              {statusLabels[taskData.status]}
            </span>
          </div>
          <p className="text-gray-500 mt-2">任务 ID: {taskData.task_id}</p>
        </div>

        {/* Progress */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">进度</h2>
          <div className="mb-2">
            <div className="flex justify-between text-sm text-gray-600 mb-1">
              <span>{taskData.progress.current_step}</span>
              <span>{taskData.progress.percentage}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-4">
              <div
                className="bg-blue-500 h-4 rounded-full transition-all duration-300"
                style={{ width: `${taskData.progress.percentage}%` }}
              />
            </div>
          </div>
          <p className="text-sm text-gray-600 mt-2">
            步骤 {taskData.progress.current_step_index} / {taskData.progress.total_steps}
          </p>

          {taskData.progress.details && (
            <div className="mt-4 grid grid-cols-2 gap-4 text-sm">
              <div className="bg-gray-50 p-3 rounded">
                <p className="text-gray-600">角色提取</p>
                <p className="text-lg font-semibold">{taskData.progress.details.characters_extracted}</p>
              </div>
              <div className="bg-gray-50 p-3 rounded">
                <p className="text-gray-600">角色参考图</p>
                <p className="text-lg font-semibold">{taskData.progress.details.characters_generated}</p>
              </div>
              <div className="bg-gray-50 p-3 rounded">
                <p className="text-gray-600">场景划分</p>
                <p className="text-lg font-semibold">{taskData.progress.details.scenes_divided}</p>
              </div>
              <div className="bg-gray-50 p-3 rounded">
                <p className="text-gray-600">场景生成</p>
                <p className="text-lg font-semibold">{taskData.progress.details.scenes_generated}</p>
              </div>
            </div>
          )}
        </div>

        {/* Result */}
        {taskData.result && (
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4">生成结果</h2>
            <div className="grid grid-cols-2 gap-4 mb-6">
              <div>
                <p className="text-gray-600">角色数量</p>
                <p className="text-2xl font-bold text-blue-600">{taskData.result.character_count}</p>
              </div>
              <div>
                <p className="text-gray-600">场景数量</p>
                <p className="text-2xl font-bold text-green-600">{taskData.result.scene_count}</p>
              </div>
            </div>

            {/* Characters */}
            {taskData.result.characters.length > 0 && (
              <div className="mb-6">
                <h3 className="font-semibold mb-3">角色</h3>
                <div className="grid grid-cols-3 gap-4">
                  {taskData.result.characters.map((char) => (
                    <div key={char.id} className="border rounded-lg p-3">
                      {char.reference_image_url && (
                        <img
                          src={char.reference_image_url}
                          alt={char.name}
                          className="w-full h-32 object-cover rounded mb-2"
                        />
                      )}
                      <p className="font-medium text-center">{char.name}</p>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Scenes */}
            {taskData.result.scenes.length > 0 && (
              <div>
                <h3 className="font-semibold mb-3">场景</h3>
                <div className="space-y-4">
                  {taskData.result.scenes.map((scene) => (
                    <div key={scene.id} className="border rounded-lg p-4">
                      <p className="font-medium mb-2">第 {scene.sequence_num} 幕</p>
                      <p className="text-gray-600 text-sm mb-2">{scene.description}</p>
                      {scene.image_url && (
                        <img
                          src={scene.image_url}
                          alt={`Scene ${scene.sequence_num}`}
                          className="w-full h-48 object-cover rounded"
                        />
                      )}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        {/* Error */}
        {taskData.error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-6 mb-6">
            <h2 className="text-xl font-semibold text-red-800 mb-2">错误信息</h2>
            <p className="text-red-600">{taskData.error.message}</p>
            <p className="text-sm text-red-500 mt-1">错误代码: {taskData.error.code}</p>
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-4">
          <button
            onClick={() => navigate('/tasks')}
            className="px-6 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700"
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
              className="px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
            >
              取消任务
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default TaskDetailPage;
