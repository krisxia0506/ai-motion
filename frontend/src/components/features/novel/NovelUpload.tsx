import React, { useState, useCallback } from 'react';
import { Button, ProgressBar, ErrorMessage } from '../../common';
import type { Novel } from '../../../types';
import './NovelUpload.css';

interface NovelUploadProps {
  onUploadSuccess?: (novel: Novel) => void;
  onUploadError?: (error: Error) => void;
}

export const NovelUpload: React.FC<NovelUploadProps> = ({
  onUploadSuccess,
  onUploadError,
}) => {
  const [file, setFile] = useState<File | null>(null);
  const [title, setTitle] = useState('');
  const [author, setAuthor] = useState('');
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const validateFile = (file: File): string | null => {
    const maxSize = 10 * 1024 * 1024;
    const allowedTypes = ['.txt', '.doc', '.docx'];
    const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase();

    if (file.size > maxSize) {
      return 'File size must be less than 10MB';
    }

    if (!allowedTypes.includes(fileExtension)) {
      return 'Only .txt, .doc, and .docx files are allowed';
    }

    return null;
  };

  const handleFileSelect = (selectedFile: File) => {
    const validationError = validateFile(selectedFile);
    if (validationError) {
      setError(validationError);
      return;
    }

    setFile(selectedFile);
    setError(null);
    if (!title) {
      setTitle(selectedFile.name.replace(/\.[^/.]+$/, ''));
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
      handleFileSelect(droppedFile);
    }
  }, []);

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      handleFileSelect(selectedFile);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!file || !title) {
      setError('Please select a file and provide a title');
      return;
    }

    setIsUploading(true);
    setError(null);
    setUploadProgress(0);

    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('title', title);
      formData.append('author', author);

      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => Math.min(prev + 10, 90));
      }, 200);

      const response = await fetch('/api/v1/novels/upload', {
        method: 'POST',
        body: formData,
      });

      clearInterval(progressInterval);

      if (!response.ok) {
        throw new Error('Upload failed');
      }

      const data = await response.json();
      setUploadProgress(100);

      setTimeout(() => {
        onUploadSuccess?.(data.novel);
        setFile(null);
        setTitle('');
        setAuthor('');
        setUploadProgress(0);
        setIsUploading(false);
      }, 500);
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Upload failed');
      setError(error.message);
      onUploadError?.(error);
      setIsUploading(false);
      setUploadProgress(0);
    }
  };

  return (
    <div className="novel-upload">
      <h2 className="novel-upload-title">Upload Novel</h2>

      <form onSubmit={handleSubmit} className="novel-upload-form">
        <div
          className={`novel-upload-dropzone ${isDragging ? 'dragging' : ''} ${
            file ? 'has-file' : ''
          }`}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          {!file ? (
            <>
              <div className="dropzone-icon">📄</div>
              <p className="dropzone-text">
                Drag and drop your novel file here, or
              </p>
              <label htmlFor="file-input" className="dropzone-browse-btn">
                Browse Files
              </label>
              <input
                id="file-input"
                type="file"
                accept=".txt,.doc,.docx"
                onChange={handleFileInputChange}
                className="dropzone-file-input"
              />
              <p className="dropzone-hint">
                Supported formats: TXT, DOC, DOCX (Max 10MB)
              </p>
            </>
          ) : (
            <div className="dropzone-file-info">
              <div className="file-info-icon">📄</div>
              <div className="file-info-details">
                <p className="file-info-name">{file.name}</p>
                <p className="file-info-size">
                  {(file.size / 1024 / 1024).toFixed(2)} MB
                </p>
              </div>
              <button
                type="button"
                onClick={() => setFile(null)}
                className="file-info-remove-btn"
              >
                ✕
              </button>
            </div>
          )}
        </div>

        {error && <ErrorMessage message={error} onRetry={() => setError(null)} />}

        <div className="novel-upload-fields">
          <div className="form-field">
            <label htmlFor="title" className="form-label">
              Title *
            </label>
            <input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Enter novel title"
              required
              disabled={isUploading}
              className="form-input"
            />
          </div>

          <div className="form-field">
            <label htmlFor="author" className="form-label">
              Author
            </label>
            <input
              id="author"
              type="text"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              placeholder="Enter author name (optional)"
              disabled={isUploading}
              className="form-input"
            />
          </div>
        </div>

        {isUploading && (
          <div className="upload-progress">
            <ProgressBar
              value={uploadProgress}
              showLabel
              label={`Uploading... ${uploadProgress}%`}
            />
          </div>
        )}

        <div className="novel-upload-actions">
          <Button
            type="submit"
            disabled={!file || !title || isUploading}
            className="upload-submit-btn"
          >
            {isUploading ? 'Uploading...' : 'Upload Novel'}
          </Button>
        </div>
      </form>
    </div>
  );
};
