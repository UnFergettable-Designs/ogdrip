export interface OpenGraphRequest {
    url?: string;
    title?: string;
    description?: string;
    type?: 'website' | 'article' | 'product' | 'profile';
    site?: string;
    targetUrl?: string;
    width?: number;
    height?: number;
    twitterCard?: 'summary' | 'summary_large_image';
    debug?: boolean;
    verbose?: boolean;
    selector?: string;
    wait?: number;
    quality?: number;
}
export interface OpenGraphResponse {
    success: boolean;
    message: string;
    imageUrl?: string;
    metaTagsUrl?: string;
    id?: string;
    error?: string;
}
