import React, { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Input, ProgressBar } from '../../common';
import { useNovelStore } from '../../../store/novelStore';
import { novelApi } from '../../../services/novelApi';
import type { UploadNovelRequest } from '../../../types';
import './NovelUpload.css';

export const NovelUpload: React.FC = () => {
  const navigate = useNavigate();
  const addNovel = useNovelStore((state) => state.addNovel);
  
  const [formData, setFormData] = useState({
    title: '',
    author: '',
    content: '',
  });
  
  const [file, setFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    setError(null);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      validateAndSetFile(selectedFile);
    }
  };

  const validateAndSetFile = (selectedFile: File) => {
    const maxSize = 10 * 1024 * 1024;
    const allowedTypes = ['text/plain', 'application/epub+zip'];
    
    if (!allowedTypes.includes(selectedFile.type) && !selectedFile.name.endsWith('.txt') && !selectedFile.name.endsWith('.epub')) {
      setError('Only TXT and EPUB files are allowed');
      return;
    }
    
    if (selectedFile.size > maxSize) {
      setError('File size must be less than 10MB');
      return;
    }
    
    setFile(selectedFile);
    setError(null);
    
    if (selectedFile.type === 'text/plain' || selectedFile.name.endsWith('.txt')) {
      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target?.result as string;
        setFormData((prev) => ({ ...prev, content }));
      };
      reader.readAsText(selectedFile);
    }
  };

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) {
      validateAndSetFile(droppedFile);
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.title.trim()) {
      setError('Title is required');
      return;
    }
    
    if (!formData.author.trim()) {
      setError('Author is required');
      return;
    }
    
    if (!formData.content.trim() && !file) {
      setError('Please provide novel content or upload a file');
      return;
    }

    setIsUploading(true);
    setError(null);
    setUploadProgress(0);

    try {
      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => Math.min(prev + 10, 90));
      }, 200);

      const request: UploadNovelRequest = {
        title: formData.title,
        author: formData.author,
        content: formData.content,
      };

      const novel = await novelApi.uploadNovel(request);
      
      clearInterval(progressInterval);
      setUploadProgress(100);
      
      addNovel(novel);
      
      setTimeout(() => {
        navigate(`/novels/${novel.id}`);
      }, 500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to upload novel');
      setUploadProgress(0);
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="novel-upload">
      <h2 className="novel-upload-title">Upload Novel</h2>
      
      <form onSubmit={handleSubmit} className="novel-upload-form">
        <Input
          label="Title"
          name="title"
          value={formData.title}
          onChange={handleInputChange}
          placeholder="Enter novel title"
          required
          fullWidth
          disabled={isUploading}
        />

        <Input
          label="Author"
          name="author"
          value={formData.author}
          onChange={handleInputChange}
          placeholder="Enter author name"
          required
          fullWidth
          disabled={isUploading}
        />

        <div className="novel-upload-file-section">
          <label className="novel-upload-label">Upload File (Optional)</label>
          
          <div
            className={`novel-upload-dropzone ${isDragging ? 'novel-upload-dropzone-dragging' : ''} ${file ? 'novel-upload-dropzone-has-file' : ''}`}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
          >
            <input
              type="file"
              accept=".txt,.epub"
              onChange={handleFileChange}
              className="novel-upload-file-input"
              id="file-input"
              disabled={isUploading}
            />
            
            {file ? (
              <div className="novel-upload-file-info">
                <svg
                  className="novel-upload-file-icon"
                  width="48"
                  height="48"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
                  <polyline points="13 2 13 9 20 9" />
                </svg>
                <div className="novel-upload-file-details">
                  <p className="novel-upload-file-name">{file.name}</p>
                  <p className="novel-upload-file-size">
                    {(file.size / 1024).toFixed(2)} KB
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => setFile(null)}
                  className="novel-upload-file-remove"
                  disabled={isUploading}
                >
                  Remove
                </button>
              </div>
            ) : (
              <label htmlFor="file-input" className="novel-upload-dropzone-content">
                <svg
                  className="novel-upload-upload-icon"
                  width="48"
                  height="48"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                  <polyline points="17 8 12 3 7 8" />
                  <line x1="12" y1="3" x2="12" y2="15" />
                </svg>
                <p className="novel-upload-dropzone-text">
                  {isDragging
                    ? 'Drop file here'
                    : 'Drag and drop a file here, or click to select'}
                </p>
                <p className="novel-upload-dropzone-hint">
                  Supports TXT and EPUB files up to 10MB
                </p>
              </label>
            )}
          </div>
        </div>

        <div className="novel-upload-textarea-wrapper">
          <label htmlFor="content" className="novel-upload-label">
            Novel Content {!file && <span className="novel-upload-required">*</span>}
          </label>
          <textarea
            id="content"
            name="content"
            value={formData.content}
            onChange={handleInputChange}
            placeholder="Paste novel content here..."
            className="novel-upload-textarea"
            rows={10}
            disabled={isUploading}
          />
        </div>

        {error && <div className="novel-upload-error">{error}</div>}

        {isUploading && (
          <ProgressBar
            value={uploadProgress}
            label="Uploading novel..."
            showPercentage
          />
        )}

        <div className="novel-upload-actions">
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate('/novels')}
            disabled={isUploading}
          >
            Cancel
          </Button>
          <Button type="submit" loading={isUploading}>
            {isUploading ? 'Uploading...' : 'Upload Novel'}
          </Button>
        </div>
      </form>
    </div>
  );
};
