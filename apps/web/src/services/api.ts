import axios from "axios";

const BACKEND_BASE_URL = process.env.NEXT_PUBLIC_BACKEND_BASE_URL || "http://localhost";
const BACKEND_PORT = process.env.NEXT_PUBLIC_BACKEND_PORT || "5000";
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || `${BACKEND_BASE_URL}:${BACKEND_PORT}/api`;

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(
  (config) => {
    if (typeof window !== "undefined") {
      const token = localStorage.getItem("token");
      if (token) {
        config.headers["Authorization"] = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default api;

export const mediaService = {
  getMediaDetails: async (source: string, externalId: string, itemType: string) => {
    return api.get(`/media/details`, {
      params: { source, external_id: externalId, item_type: itemType },
    });
  },

  getTrendingMedia: async (limit = 10) => {
    return api.get(`/media/trending?limit=${limit}`);
  },
};