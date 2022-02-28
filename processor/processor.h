#ifndef PROCESSOR_H
#define PROCESSOR_H
#include <darknet.h>
typedef struct BoundingBox {
    int typeID;
    float height;
    float width;
    float posx;
    float posy;
} BoundingBox;

typedef struct {
    char *cfg_file;
    char *weight_file;
    char *names_file;
    char **labels;
    size_t classes;
    float thresh;
    float hier_thresh;
    network *net;
    detection *detections;
} ProcessServer;

ProcessServer *newProcessServer();
int loadModel(ProcessServer *self);
// void detect(ProcessServer *self, char *filepath,void *cgoBoundingboxes);
void detectFuck(void *cgoBoundingboxes);
void getBoundingBox(ProcessServer *self, int num_boxes, void *cgoBoundingboxes);
int runDetection(ProcessServer *self, char *filepath);
// void Hello();
#endif