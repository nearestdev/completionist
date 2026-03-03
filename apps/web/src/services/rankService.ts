import { UserRank } from "@/types/rank";
import api from "./api";

export const getUserRank = async (): Promise<UserRank> => {
    const response = await api.get<UserRank>("/me/rank");
    return response.data;
};
