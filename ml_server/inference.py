import torch
from PIL import Image
import numpy as np
from torchvision import transforms
from models.UANet import UANet_VGG
import matplotlib.pyplot as plt
import os
import rasterio  # 添加 rasterio 库来处理 TIF 图像
import argparse

class BuildingPredictor:
    def __init__(self, checkpoint_path):
        self.device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
        self.model = UANet_VGG(channel=32, num_classes=2)
        
        checkpoint = torch.load(checkpoint_path, map_location=self.device)
        state_dict = checkpoint['state_dict']
        
        new_state_dict = {}
        for k, v in state_dict.items():
            if k.startswith('net.'):
                new_state_dict[k[4:]] = v
            else:
                new_state_dict[k] = v
                
        self.model.load_state_dict(new_state_dict)
        self.model.to(self.device)
        self.model.eval()
        
        self.transform = transforms.Compose([
            transforms.Resize((512, 512)),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.485, 0.456, 0.406],
                              std=[0.229, 0.224, 0.225])
        ])
    
    def load_image(self, image_path):
        """
        加载图像，支持 TIF 和普通图像格式
        """
        if image_path.lower().endswith('.tif') or image_path.lower().endswith('.tiff'):
            # 使用 rasterio 读取 TIF 图像
            with rasterio.open(image_path) as src:
                # 读取所有波段
                image = src.read()
                # 转换为 RGB 格式 (如果需要)
                if image.shape[0] > 3:
                    image = image[:3]  # 只使用前三个波段
                # 转换为 PIL Image 格式
                image = np.transpose(image, (1, 2, 0))  # (C,H,W) -> (H,W,C)
                # 标准化到 0-255 范围
                image = ((image - image.min()) / (image.max() - image.min()) * 255).astype(np.uint8)
                return Image.fromarray(image).convert('RGB')
        else:
            # 普通图像格式使用 PIL
            return Image.open(image_path).convert('RGB')

    def predict(self, image_path, save_path=None):
        """
        预测单张图片
        """
        # 加载图片
        image = self.load_image(image_path)
        original_size = image.size
        
        # 预处理
        input_tensor = self.transform(image)
        input_batch = input_tensor.unsqueeze(0).to(self.device)
        
        # 预测
        with torch.no_grad():
            output = self.model(input_batch)
            if isinstance(output, tuple):
                output = output[-1]
            
        # 后处理
        pred = torch.softmax(output, dim=1)
        pred = pred.argmax(dim=1).cpu().numpy()[0]
        
        # 将预测结果转换为二值图像
        mask = (pred * 255).astype(np.uint8)
        
        # 调整回原始大小
        mask_image = Image.fromarray(mask).resize(original_size, Image.NEAREST)
        
        # 保存成 PNG
        if save_path:
            # 将传入的 save_path 扩展改为 _mask.png
            mask_png_path = os.path.splitext(save_path)[0] + '_mask.png'
            mask_image.save(mask_png_path, 'PNG')
        
        return mask_image

if __name__ == "__main__":
    # 创建命令行参数解析器
    parser = argparse.ArgumentParser(description='建筑物检测预测程序')
    parser.add_argument('--input', '-i', type=str, required=True, help='输入图片路径')
    parser.add_argument('--output', '-o', type=str, default='result.png', help='输出结果路径, 不带后缀_mask')
    parser.add_argument('--model', '-m', type=str, default='models/UANet_VGG.ckpt', help='模型检查点路径')
    
    args = parser.parse_args()
    
    # 使用命令行参数
    predictor = BuildingPredictor(args.model)
    result = predictor.predict(args.input, args.output) 