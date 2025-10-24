import React from 'react';
import './Select.css';

export interface SelectOption {
  value: string;
  label: string;
  disabled?: boolean;
}

interface SelectProps extends Omit<React.SelectHTMLAttributes<HTMLSelectElement>, 'onChange'> {
  label?: string;
  error?: string;
  helperText?: string;
  fullWidth?: boolean;
  options: SelectOption[];
  placeholder?: string;
  onChange?: (value: string) => void;
}

export const Select: React.FC<SelectProps> = ({
  label,
  error,
  helperText,
  fullWidth = false,
  options,
  placeholder,
  onChange,
  className = '',
  id,
  value,
  ...props
}) => {
  const selectId = id || `select-${Math.random().toString(36).substr(2, 9)}`;
  const classes = ['select-wrapper', fullWidth && 'select-full-width', className]
    .filter(Boolean)
    .join(' ');

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    onChange?.(e.target.value);
  };

  return (
    <div className={classes}>
      {label && (
        <label htmlFor={selectId} className="select-label">
          {label}
          {props.required && <span className="select-required"> *</span>}
        </label>
      )}
      <select
        id={selectId}
        className={`select ${error ? 'select-error' : ''}`}
        onChange={handleChange}
        value={value}
        {...props}
      >
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}
        {options.map((option) => (
          <option
            key={option.value}
            value={option.value}
            disabled={option.disabled}
          >
            {option.label}
          </option>
        ))}
      </select>
      {error && <span className="select-error-text">{error}</span>}
      {helperText && !error && (
        <span className="select-helper-text">{helperText}</span>
      )}
    </div>
  );
};
