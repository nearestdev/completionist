import api from "./api";
import { Franchise, FranchiseWithItems } from "@/types/franchise";

const franchiseService = {
  list: async (limit: number = 20, offset: number = 0, category?: string): Promise<Franchise[]> => {
    let url = `/franchises?limit=${limit}&offset=${offset}`;
    if (category) url += `&category=${category}`;
    const response = await api.get<Franchise[]>(url);
    return response.data;
  },
  getById: async (id: number): Promise<FranchiseWithItems> => {
    const response = await api.get<FranchiseWithItems>(`/franchises/${id}`);
    return response.data;
  },
};
export default franchiseService;
