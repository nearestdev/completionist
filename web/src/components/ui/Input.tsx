import React, { forwardRef } from "react";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className, ...props }, ref) => {
    return (
      <input
        ref={ref}
        className={`mt-1 block w-full px-3 py-2 border rounded-md shadow-sm 
                   bg-white dark:bg-gray-800 
                   border-gray-300 dark:border-gray-600 
                   text-gray-900 dark:text-gray-200 
                   placeholder-gray-400 dark:placeholder-gray-500 
                   focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 
                   ${className}`}
        {...props}
      />
    );
  }
);

Input.displayName = "Input";
export default Input;