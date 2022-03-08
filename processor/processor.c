#include <stdio.h>
#include <darknet.h>
#include "processor.h"

// rm processor.o && gcc -c processor.c -L. -ldarknet -Wl,-rpath,$PWD/libdarknet.so




// static char *cfg_file = "oiltank-playground-yolov3.cfg";
// static char *weight_file = "oiltank-playground-yolov3_900.weights";
// static char *names_file = "oiltank-playground-obj.names";

ProcessServer *newProcessServer() {
    ProcessServer *ps = malloc(sizeof(ProcessServer));
    ps->cfg_file = "processor/oiltank-yolov2.cfg";
    ps->weight_file = "processor/oiltank-yolov2_900.weights";
    ps->names_file = "processor/oiltank-yolov2.names";
    ps->labels = get_labels(ps->names_file);

    ps->thresh = 0.5;
    ps->hier_thresh = 0.5;


    ps->classes = 0;
    ps->detections = NULL;
    // char **labels = get_labels(names_file);
    while (ps->labels[ps->classes] != NULL) {
        ps->classes++;
    }
    return ps;
}

int loadModel(ProcessServer *self) {
    network *net = load_network(self->cfg_file, self->weight_file, 0);
    self->net = net;
    set_batch_network(net, 1);
    return 0;
}

int runDetection(ProcessServer *self, char *filepath) {
    //return 0;
    printf("running detection on : %s", filepath);
    image im = load_image_color(filepath, 0, 0);
    
    // And scale it to the parameters define din *.cfg file.
    image sized = letterbox_image(im, self->net->w, self->net->h);

    
    float *frame_data = sized.data;

    // Do prediction.
    double time = what_time_is_it_now();
    network_predict(self->net, frame_data);

    printf("'%s' predicted in %lf seconds\n", filepath, 
        (what_time_is_it_now() - time));

    // Get number fo predicted classes (objects).
    int num_boxes = 0;
    self->detections = get_network_boxes(self->net, 
        im.w, im.h, self->thresh, self->hier_thresh, NULL, 1, &num_boxes);

    printf("Detected %d object, class %d", num_boxes, self->detections->classes);
    free_image(im);
    free_image(sized);
    return num_boxes;
}

void getBoundingBox(ProcessServer *self, int num_boxes, void *cgoBoundingboxes) {
    BoundingBox *boundingboxes = (BoundingBox *)cgoBoundingboxes;
    

    // -----------------------------------------------------------------------------------------------------------------
    // Print results.
    // -----------------------------------------------------------------------------------------------------------------

    // Iterate over predicted classes and print information.
	// Overflow when num_boxes >= 256
	// should use int

    for (int i = 0; i < num_boxes; ++i) {
        for (int j = 0; j < self->classes; ++j) {
            if (self->detections[i].prob[j] > self->thresh) {
                printf("Num of Boxes: %d, current i : %d", num_boxes, i);
                printf("%s %d ", self->labels[j],  (int16_t) (self->detections[i].prob[j] * 100));
                // BoundingBox boundingbox = boundingboxes[i];
                printf("%f, %f, %f, %f \n",
                    self->detections[i].bbox.h,self->detections[i].bbox.w, 
                    self->detections[i].bbox.x, self->detections[i].bbox.y);
                boundingboxes[i].typeID = j;
                boundingboxes[i].height = self->detections[i].bbox.h;
                boundingboxes[i].width = self->detections[i].bbox.w;
                boundingboxes[i].posx = self->detections[i].bbox.x;
                boundingboxes[i].posy = self->detections[i].bbox.y;
                // boundingboxes[i] = boundingbox;
            }
        }
    }


    
    // -----------------------------------------------------------------------------------------------------------------
    // Free resources.
    // -----------------------------------------------------------------------------------------------------------------

    free_detections(self->detections, num_boxes);
    return;


}
