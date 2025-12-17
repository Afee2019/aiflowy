import audioIcon from '#/assets/ai/resource/audio-icon.png';
import docIcon from '#/assets/ai/resource/doc-icon.png';
import otherIcon from '#/assets/ai/resource/other-icon.png';
import videoIcon from '#/assets/ai/resource/video-icon.png';

export function getSrc(item: any) {
  switch (item.resourceType) {
    case 0: {
      return item.resourceUrl;
    }
    case 1: {
      return audioIcon;
    }
    case 2: {
      return videoIcon;
    }
    case 3: {
      return docIcon;
    }
    default: {
      return otherIcon;
    }
  }
}

export function getResourceTypeColor(item: any) {
  switch (item.resourceType) {
    case 0: {
      return '#0066FF';
    }
    case 1: {
      return '#FFA200';
    }
    case 2: {
      return '#5600FF';
    }
    case 3: {
      return '#0099CC';
    }
    default: {
      return '#757575';
    }
  }
}
export function getResourceOriginColor(item: any) {
  switch (item.origin) {
    case 0: {
      return '#039e90';
    }
    case 1: {
      return '#0066FF';
    }
    default: {
      return '#757575';
    }
  }
}
